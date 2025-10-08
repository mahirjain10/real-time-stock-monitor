import 'package:flutter/material.dart';
import 'package:frontend/screens/auth/presentation/login_page.dart';
import 'package:frontend/screens/auth/presentation/signup_page.dart';
import 'package:frontend/main.dart';

class SignUpAndLoginText extends StatelessWidget {
  final VoidCallback? onTap;
  final bool isLogin;
  const SignUpAndLoginText({Key? key, required this.isLogin, this.onTap})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    String text = isLogin ? "Don't have an account? " : "Have an account? ";
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Text(text, style: const TextStyle(fontSize: 15)),
        GestureDetector(
          onTap: onTap ??
              () {
                print("Sign Up Clicked!");
                isLogin
                    ? Navigator.push(
                        context,
                        MaterialPageRoute(
                            builder: (context) => const SignUpPage()),
                      )
                    : Navigator.push(
                        context,
                        MaterialPageRoute(
                            builder: (context) => LoginPage(dio: dio)),
                      );
              },
          child: Text(
            isLogin ? "Sign Up" : "Login",
            style: const TextStyle(
              color: Color(0xff2563eb),
              fontSize: 15,
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
      ],
    );
  }
}
