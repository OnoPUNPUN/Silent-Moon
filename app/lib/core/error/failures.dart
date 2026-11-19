import 'package:equatable/equatable.dart';

abstract class Failure extends Equatable {
  final String message;

  const new(this.message);

  @override
  List<Object?> get props => [message];
}

class ServerFailure extends Failure {
  const new(super.message);
}

class CacheFailure extends Failure {
  const new(super.message);
}
