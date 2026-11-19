import 'dart:convert';

import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:logger/logger.dart';

/// A custom Dio interceptor that formats API requests, responses, and errors
/// into clean, organized, and colorful boxed logs using the [Logger] package.
class DioLoggerInterceptor extends Interceptor {
  /// Whether to log HTTP requests.
  final bool request;

  /// Whether to log request headers. (Kept false by default to prevent console clutter)
  final bool requestHeader;

  /// Whether to log request body data.
  final bool requestBody;

  /// Whether to log response headers. (Kept false by default to remove Cloudflare/server header spam)
  final bool responseHeader;

  /// Whether to log response body data.
  final bool responseBody;

  /// Whether to log error responses.
  final bool error;

  /// Whether to track and display API response time in milliseconds.
  final bool logDuration;

  /// Whether logging is enabled. Defaults to [kDebugMode] so production builds stay clean and safe.
  final bool enabled;

  /// Max characters to display for request/response bodies before truncating.
  final int maxContentLength;

  late final Logger _logger;

  DioLoggerInterceptor({
    this.request = true,
    this.requestHeader = false,
    this.requestBody = true,
    this.responseHeader = false,
    this.responseBody = true,
    this.error = true,
    this.logDuration = true,
    this.enabled = kDebugMode,
    this.maxContentLength = 10000,
    Logger? customLogger,
  }) {
    _logger =
        customLogger ??
        Logger(
          printer: PrettyPrinter(
            methodCount: 0,
            errorMethodCount: 0,
            lineLength: 85,
            colors: true,
            printEmojis: true,
            dateTimeFormat: DateTimeFormat.none,
          ),
        );
  }

  static const _startTimeKey = 'dio_logger_start_time';

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    if (enabled) {
      if (logDuration) {
        options.extra[_startTimeKey] = DateTime.now().millisecondsSinceEpoch;
      }

      if (request) {
        final buffer = StringBuffer();
        buffer.writeln('🚀 [REQUEST] ${options.method} ${options.uri}');

        if (options.queryParameters.isNotEmpty) {
          buffer.writeln('Query: ${_prettyJson(options.queryParameters)}');
        }

        if (requestHeader && options.headers.isNotEmpty) {
          buffer.writeln('Headers: ${_prettyJson(options.headers)}');
        }

        if (requestBody && options.data != null) {
          buffer.writeln('Body:\n${_prettyJson(options.data)}');
        }

        _logger.i(_sanitize(buffer.toString()));
      }
    }
    super.onRequest(options, handler);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    if (enabled) {
      final duration = _calculateDuration(response.requestOptions);
      final durationStr = duration != null ? ' ($duration ms)' : '';
      final statusCode = response.statusCode ?? 0;
      final statusMessage = response.statusMessage ?? '';

      final buffer = StringBuffer();
      buffer.writeln(
        '✅ [RESPONSE $statusCode ${statusMessage.isNotEmpty ? statusMessage : 'OK'}]$durationStr ${response.requestOptions.method} ${response.requestOptions.uri}',
      );

      if (responseHeader && response.headers.map.isNotEmpty) {
        buffer.writeln('Headers: ${_prettyJson(response.headers.map)}');
      }

      if (responseBody && response.data != null) {
        buffer.writeln('Data:\n${_prettyJson(response.data)}');
      }

      final message = _sanitize(buffer.toString());
      if (statusCode >= 200 && statusCode < 300) {
        _logger.i(message);
      } else if (statusCode >= 300 && statusCode < 400) {
        _logger.w(message);
      } else {
        _logger.e(message);
      }
    }
    super.onResponse(response, handler);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    if (enabled && error) {
      final duration = _calculateDuration(err.requestOptions);
      final durationStr = duration != null ? ' ($duration ms)' : '';
      final statusCode = err.response?.statusCode;
      final statusMessage = err.response?.statusMessage;
      final statusStr = statusCode != null
          ? ' $statusCode${statusMessage != null ? ' $statusMessage' : ''}'
          : '';

      final buffer = StringBuffer();
      buffer.writeln(
        '❌ [ERROR$statusStr]$durationStr ${err.requestOptions.method} ${err.requestOptions.uri}',
      );
      buffer.writeln('Type: ${err.type.name}');

      if (err.message != null && err.message!.isNotEmpty) {
        buffer.writeln('Message: ${err.message}');
      }

      if (err.response?.data != null) {
        buffer.writeln('Response Data:\n${_prettyJson(err.response?.data)}');
      }

      _logger.e(_sanitize(buffer.toString()), error: err.error);
    }
    super.onError(err, handler);
  }

  int? _calculateDuration(RequestOptions options) {
    final startTime = options.extra[_startTimeKey] as int?;
    if (startTime == null) return null;
    return DateTime.now().millisecondsSinceEpoch - startTime;
  }

  String _sanitize(String content) {
    return content.replaceAll('\r\n', '\n').trimRight();
  }

  String _prettyJson(dynamic data) {
    if (data == null) return 'null';
    try {
      String jsonStr;
      if (data is Map || data is List) {
        jsonStr = const JsonEncoder.withIndent('  ').convert(data);
      } else if (data is String) {
        final trimmed = data.trim();
        if ((trimmed.startsWith('{') && trimmed.endsWith('}')) ||
            (trimmed.startsWith('[') && trimmed.endsWith(']'))) {
          final decoded = jsonDecode(trimmed);
          jsonStr = const JsonEncoder.withIndent('  ').convert(decoded);
        } else {
          jsonStr = data;
        }
      } else if (data is FormData) {
        final fields = Map.fromEntries(data.fields);
        final files = data.files
            .map((e) => '${e.key}: ${e.value.filename} (${e.value.length} B)')
            .toList();
        jsonStr = const JsonEncoder.withIndent('  ')
            .convert({'fields': fields, 'files': files});
      } else {
        jsonStr = data.toString();
      }

      if (jsonStr.length > maxContentLength) {
        return '${jsonStr.substring(0, maxContentLength)}\n... (truncated ${jsonStr.length - maxContentLength} characters)';
      }
      return jsonStr;
    } catch (_) {
      return data.toString();
    }
  }
}
