class ServerExpection implements Exception {
  final String message;

  new({this.message = "Server Exception"});
}

class CacheException implements Exception {}
