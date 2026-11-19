import 'package:dio/dio.dart';

class ApiClinet {
  final Dio dio;

  new(this.dio);

  Future<Response> get(
    String path, {
    Map<String, dynamic>? queryParams,
    Options? options,
  }) async {
    return await dio.get(path, queryParameters: queryParams, options: options);
  }

  Future<Response> post(
    String path,
    dynamic data,
    Map<String, dynamic>? queryParams,
    Options? options,
  ) async {
    return await dio.post(
      path,
      data: data,
      queryParameters: queryParams,
      options: options,
    );
  }

  Future<Response> put(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParams,
    Options? options,
  }) async {
    return await dio.put(
      path,
      data: data,
      queryParameters: queryParams,
      options: options,
    );
  }

  Future<Response> patch(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParams,
    Options? options,
  }) async {
    return await dio.patch(
      path,
      data: data,
      queryParameters: queryParams,
      options: options,
    );
  }

  Future<Response> delete(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParams,
    Options? options,
  }) async {
    return await dio.delete(
      path,
      data: data,
      queryParameters: queryParams,
      options: options,
    );
  }
}
