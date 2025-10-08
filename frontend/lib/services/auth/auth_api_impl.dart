import 'dart:convert';
import 'package:dio/dio.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:frontend/models/response_model.dart';
import 'package:frontend/models/user_model.dart';
import 'package:frontend/services/auth/auth_api_service.dart';

class AuthApiImpl implements AuthApiService {
  final Dio dio;
  final String baseUrl = dotenv.env["BASE_URL"]!;

  AuthApiImpl(this.dio) {
    dio.options.baseUrl = baseUrl;
    dio.options.headers["Content-Type"] = "application/json";
  }

  @override
  Future<ResponseModel<LoginUserModel>> login(
      String email, String password) async {
    try {
      Uri url = Uri.parse("$baseUrl/api/auth/login");

      print("Login request: $email, $password");
      print("Login password length: ${password.length}");

      LoginUserModel loginModel =
          LoginUserModel(email: email, password: password);
      String requestBody = loginModel.toJsonString();
      print("Request body: $requestBody");
      final response = await dio.post(
        "/api/auth/login",
        data: requestBody,
      );

      print("Response status code: ${response.statusCode}");
      print("Response body: ${response.data}");
      print("Response headers: ${response.headers}");

      return ResponseModel<LoginUserModel>.fromJson(
        jsonEncode(response.data),
        response.statusCode!,
        response.data,
        (data) => LoginUserModel.fromJson(data),
      );
    } catch (e) {
      throw Exception("Error during login: $e");
    }
  }

  @override
  Future<ResponseModel<SignupUserModel>> register(
      String name, String email, String password) async {
    try {
      Uri url = Uri.parse("$baseUrl/api/auth/register");

      print("Register request: $name, $email, $password");
      print("Register password length: ${password.length}");

      SignupUserModel registerModel =
          SignupUserModel(name: name, email: email, password: password);
      String requestBody = registerModel.toJsonString();

      final response = await dio.post(
        "/api/auth/register",
        data: requestBody,
      );

      print("Response status code: ${response.statusCode}");
      print("Response body: ${response.data}");
      print("Response headers: ${response.headers}");

      return ResponseModel<SignupUserModel>.fromJson(
        jsonEncode(response.data),
        response.statusCode!,
        response.data,
        null,
      );
    } catch (e) {
      throw Exception("Error during registration: $e");
    }
  }
}
