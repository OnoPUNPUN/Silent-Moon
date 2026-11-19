import 'package:app/core/network/dio_logger_interceptor.dart';
import 'package:app/core/network/token_manager.dart';
import 'package:dio/dio.dart';

class DioProvider {
  static const baseUrl = "http://localhost:8080/api/v1";

  static Dio createDio(TokenManager tokenManager) {
    final dio = Dio(
      BaseOptions(
        baseUrl: baseUrl,
        connectTimeout: const Duration(seconds: 10),
        receiveTimeout: const Duration(seconds: 10),
        headers: {'Content-Type': 'application/json'},
      ),
    );

    dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          final isAuthRequest = options.path.startsWith('/auth/');

          if (!isAuthRequest) {
            final token = await tokenManager.getAccessToken();
            if (token != null && token.isNotEmpty) {
              options.headers['Authorization'] = 'Bearer $token';
            }
          }

          return handler.next(options);
        },
      ),
    );

    dio.interceptors.add(
      DioLoggerInterceptor(
        request: true,
        requestBody: true,
        requestHeader: false,
        responseHeader: false,
        responseBody: true,
        error: true,
      ),
    );
    return dio;
  }
}
