import 'package:flutter/material.dart';
import 'package:frontend/main.dart';
import 'package:frontend/screens/auth/controller/signup_controller.dart';
import 'package:frontend/screens/auth/presentation/signup_login_text.dart';
import 'package:frontend/screens/auth/presentation/welcome_text.dart';
import 'package:frontend/services/auth/auth_api_impl.dart';
import 'package:frontend/widgets/auth_button.dart';
import 'package:frontend/widgets/auth_input_bar.dart';
import 'package:frontend/widgets/auth_input_label.dart';
import 'package:toastification/toastification.dart';
import 'package:provider/provider.dart';

class SignUpPage extends StatefulWidget {
  const SignUpPage({super.key});

  @override
  State<SignUpPage> createState() => _SignUpPageState();
}

class _SignUpPageState extends State<SignUpPage> {
  void _showToast(BuildContext context, String message, bool isError) {
    debugPrint(message);
    toastification.show(
      context: context,
      type: isError ? ToastificationType.error : ToastificationType.success,
      style: ToastificationStyle.flat,
      title: Text(isError ? "Error" : "Success",
          style: const TextStyle(fontWeight: FontWeight.bold)),
      description: Text(message),
      alignment: Alignment.topCenter,
      autoCloseDuration: const Duration(seconds: 4),
    );
  }

  @override
  Widget build(BuildContext context) {
    final double screenHeight = MediaQuery.sizeOf(context).height;
    final double screenWidth = MediaQuery.sizeOf(context).width;

    return ChangeNotifierProvider(
      create: (context) => SignupController(authService: AuthApiImpl(dio)),
      child: Consumer<SignupController>(
        builder: (context, signupController, child) {
          return SafeArea(
            child: Scaffold(
              body: SingleChildScrollView(
                  child: Container(
                // alignment: AlignmentDirectional.centerStart,
                // color: Colors.red,
                child: SizedBox(
                  height: screenHeight,
                  width: screenWidth, // Take full screen height
                  child: Column(
                    mainAxisAlignment:
                        MainAxisAlignment.center, // Center vertically
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      /// **Welcome Text**
                      Padding(
                        padding: EdgeInsets.only(bottom: screenHeight * 0.05),
                        child: const WelcomeText(),
                      ),
                      Container(
                        width: screenWidth * 0.85,
                        child: Form(
                          key: signupController.formKey,
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              /// **Name Input Field**
                              const AuthInputLabel(label: "Name"),
                              InputField(
                                hintText: "Enter your name",
                                field: "name",
                                isVisible: false,
                                validateForm: (field, error) => signupController
                                    .validateField("name", error),
                                controller: signupController.nameController,
                                errorText: signupController.errorText,
                              ),
                              const SizedBox(height: 20),

                              /// **Email Input Field**
                              const AuthInputLabel(label: "Email"),
                              InputField(
                                hintText: "Enter your email",
                                field: "email",
                                isVisible: false,
                                validateForm: (field, error) => signupController
                                    .validateField("email", error),
                                controller: signupController.emailController,
                                errorText: signupController.errorText,
                              ),
                              const SizedBox(height: 20),

                              /// **Password Input Field**
                              const AuthInputLabel(label: "Password"),
                              InputField(
                                hintText: "Enter your password",
                                field: "password",
                                isVisible: signupController.isPasswordVisible,
                                validateForm: (field, error) => signupController
                                    .validateField("password", error),
                                onToggle:
                                    signupController.togglePasswordVisibility,
                                controller: signupController.passwordController,
                                errorText: signupController.errorText,
                              ),
                              const SizedBox(height: 30),

                              /// **Sign-Up Button**
                              AuthButton(
                                buttonText: "Login",
                                isLoading: signupController.isLoading,
                                onPressed: signupController.isLoading
                                    ? null
                                    : () async {
                                        await signupController.signup();
                                        // Show toast after login attempt
                                        Future.microtask(() {
                                          if (signupController.errorMessage !=
                                              null) {
                                            _showToast(
                                                context,
                                                signupController.errorMessage!,
                                                true);
                                          } else if (signupController
                                                  .successMessage !=
                                              null) {
                                            _showToast(
                                                context,
                                                signupController
                                                    .successMessage!,
                                                false);
                                          }
                                        });
                                      },
                              ),

                              /// **SignUp & Login Toggle**
                              const SizedBox(height: 20),
                              const SignUpAndLoginText(isLogin: false),
                            ],
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              )),
            ),
            // ),
          );
        },
      ),
    );
  }
}
