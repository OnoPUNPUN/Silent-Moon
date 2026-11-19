import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/// Manages secure persistence of authentication tokens, device identification,
/// and user details using platform-encrypted storage ([FlutterSecureStorage]).
class TokenManager {
  final FlutterSecureStorage _secureStorage;

  static const _accessTokenKey = 'ACCESS_TOKEN';
  static const _refreshTokenKey = 'REFRESH_TOKEN';
  static const _userIdKey = 'USER_ID';
  static const _deviceIdKey = 'DEVICE_ID';

  /// Creates a [TokenManager] instance.
  /// If [secureStorage] is omitted, standard encrypted options are used.
  TokenManager([FlutterSecureStorage? secureStorage])
      : _secureStorage = secureStorage ?? const FlutterSecureStorage();

  /// Saves access token, optional refresh token, user ID, and device ID simultaneously.
  Future<void> saveAuthTokens({
    required String accessToken,
    String? refreshToken,
    String? userId,
    String? deviceId,
  }) async {
    final futures = <Future<void>>[
      saveAccessToken(accessToken),
      if (refreshToken != null) saveRefreshToken(refreshToken),
      if (userId != null) saveUserId(userId),
      if (deviceId != null) saveDeviceId(deviceId),
    ];
    await Future.wait(futures);
  }

  // --- Access Token ---

  Future<void> saveAccessToken(String token) async {
    await _secureStorage.write(key: _accessTokenKey, value: token);
  }

  Future<String?> getAccessToken() async {
    return await _secureStorage.read(key: _accessTokenKey);
  }

  Future<void> deleteAccessToken() async {
    await _secureStorage.delete(key: _accessTokenKey);
  }

  Future<bool> hasAccessToken() async {
    final token = await getAccessToken();
    return token != null && token.isNotEmpty;
  }

  // --- Refresh Token ---

  Future<void> saveRefreshToken(String token) async {
    await _secureStorage.write(key: _refreshTokenKey, value: token);
  }

  Future<String?> getRefreshToken() async {
    return await _secureStorage.read(key: _refreshTokenKey);
  }

  Future<void> deleteRefreshToken() async {
    await _secureStorage.delete(key: _refreshTokenKey);
  }

  Future<bool> hasRefreshToken() async {
    final token = await getRefreshToken();
    return token != null && token.isNotEmpty;
  }

  // --- Convenience aliases for demo backward compatibility ---

  Future<void> saveToken(String token) async {
    await saveAccessToken(token);
  }

  Future<String?> getToken() async {
    return await getAccessToken();
  }

  // --- User ID ---

  Future<void> saveUserId(String userId) async {
    await _secureStorage.write(key: _userIdKey, value: userId);
  }

  Future<String?> getUserId() async {
    return await _secureStorage.read(key: _userIdKey);
  }

  Future<void> deleteUserId() async {
    await _secureStorage.delete(key: _userIdKey);
  }

  // --- Device ID ---

  Future<void> saveDeviceId(String deviceId) async {
    await _secureStorage.write(key: _deviceIdKey, value: deviceId);
  }

  Future<String?> getDeviceId() async {
    return await _secureStorage.read(key: _deviceIdKey);
  }

  Future<void> deleteDeviceId() async {
    await _secureStorage.delete(key: _deviceIdKey);
  }

  // --- Session Cleanup ---

  /// Clears authentication tokens and user information on logout or expired session,
  /// but preserves [deviceId] to retain device hardware binding.
  Future<void> clearAuth() async {
    await Future.wait([
      deleteAccessToken(),
      deleteRefreshToken(),
      deleteUserId(),
    ]);
  }

  /// Completely wipes all stored keys and values in secure storage.
  Future<void> clearAll() async {
    await _secureStorage.deleteAll();
  }
}
