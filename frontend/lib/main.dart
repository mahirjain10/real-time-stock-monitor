import 'dart:io';
import "package:flutter_dotenv/flutter_dotenv.dart";
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import 'package:frontend/firebase_options.dart';
import 'package:frontend/providers/auth_provider.dart';
import 'package:frontend/screens/create_alert/create_alert_page.dart';
import 'package:frontend/services/auth/auth_api_impl.dart';
import 'package:frontend/screens/auth/presentation/login_page.dart';
import 'package:frontend/screens/dashboard/dashboard_page.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';
import 'package:toastification/toastification.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:dio/dio.dart';
import 'package:dio_cookie_manager/dio_cookie_manager.dart';
import 'package:cookie_jar/cookie_jar.dart';

late final Dio dio;
late final CookieJar cookieJar;

// Background Message Handler
Future<void> _firebaseMessagingBackgroundHandler(RemoteMessage message) async {
  print("Background Message Received: ${message.notification?.title}");
}

void requestPermission() async {
  FirebaseMessaging messaging = FirebaseMessaging.instance;
  NotificationSettings settings = await messaging.requestPermission(
    alert: true,
    badge: true,
    sound: true,
  );

  if (settings.authorizationStatus == AuthorizationStatus.authorized) {
    print('User granted permission');
  } else {
    print('User denied permission');
  }
}

void getFCMToken() async {
  String? token = await FirebaseMessaging.instance.getToken();
  print("FCM Token: $token");
}

void setupFCMListeners() {
  FirebaseMessaging.onMessage.listen((RemoteMessage message) {
    print("Foreground Message: ${message.notification?.title}");
  });

  FirebaseMessaging.onMessageOpenedApp.listen((RemoteMessage message) {
    print("Notification Clicked: ${message.notification?.title}");
  });
}

Dio setupDio() {
  final dioInstance = Dio();
  cookieJar = CookieJar();

  // Add Cookie Manager to Dio
  dioInstance.interceptors.add(CookieManager(cookieJar));

  // Set base options
  dioInstance.options.baseUrl =
      dotenv.env['API_BASE_URL'] ?? 'http://localhost:8000/';
  dioInstance.options.connectTimeout = const Duration(seconds: 10);
  dioInstance.options.receiveTimeout = const Duration(seconds: 10);
  return dioInstance;
}

void main() async {
  await dotenv.load(fileName: ".env"); // Load environment variables
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp(
    options: DefaultFirebaseOptions.currentPlatform,
  );

  // Initialize Dio and Cookie Manager
  dio = setupDio();

  // Create notification channel for Android
  if (Platform.isAndroid) {
    const AndroidNotificationChannel channel = AndroidNotificationChannel(
      'default_channel',
      'Default Channel',
      description: 'Default notification channel',
      importance: Importance.high,
    );

    final FlutterLocalNotificationsPlugin flutterLocalNotificationsPlugin =
        FlutterLocalNotificationsPlugin();

    await flutterLocalNotificationsPlugin
        .resolvePlatformSpecificImplementation<
            AndroidFlutterLocalNotificationsPlugin>()
        ?.createNotificationChannel(channel);
  }

  // Setup Firebase Messaging
  FirebaseMessaging.onBackgroundMessage(_firebaseMessagingBackgroundHandler);
  requestPermission();
  getFCMToken();
  setupFCMListeners();

  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(
          create: (_) => AuthProvider(authService: AuthApiImpl(dio)),
        ),
      ],
      child: const MyApp(),
    ),
  );
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return ToastificationWrapper(
      child: MaterialApp(
        title: 'Stock Alert App',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
          textTheme: GoogleFonts.interTextTheme().copyWith(
            headlineLarge: const TextStyle(
              fontSize: 25,
              fontWeight: FontWeight.w800,
            ),
            bodyMedium: const TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.w300,
            ),
          ),
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue),
          useMaterial3: true,
        ),
        initialRoute: '/login',
        routes: {
          '/': (context) => LoginPage(dio: dio),
          '/login': (context) => LoginPage(dio: dio),
          '/dashboard': (context) => const DashboardPage(),
          '/create-alert': (context) => const CreateAlertPage(),
        },
      ),
    );
  }
}
