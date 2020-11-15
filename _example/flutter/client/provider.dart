// GENERATED CODE - DO NOT MODIFY BY HAND
import 'package:dio/dio.dart';
import 'package:flutter/widgets.dart';
import 'package:provider/provider.dart';
import 'package:provider/single_child_widget.dart';

import 'owner.dart';
import 'pet.dart';
import 'skip_generation_model.dart';

class ClientProvider extends SingleChildStatelessWidget {
  ClientProvider({
    Key key,
    @required this.dio,
    this.child,
  })  : assert(dio != null),
        super(key: key, child: child);

  final Dio dio;
  final Widget child;

  @override
  Widget buildWithChild(BuildContext context, Widget child) {
    return MultiProvider(
      providers: [
        Provider<OwnerClient>(
          create: (_) => OwnerClient(dio: dio),
        ),
        Provider<PetClient>(
          create: (_) => PetClient(dio: dio),
        ),
      ],
      child: child,
    );
  }
}
