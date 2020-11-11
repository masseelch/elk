// GENERATED CODE - DO NOT MODIFY BY HAND
import 'package:dio/dio.dart';
import 'package:flutter/widgets.dart';
import 'package:provider/provider.dart';
import 'package:provider/single_child_widget.dart';

import 'owner.dart';
import 'pet.dart';

typedef PrefixFn = String Function(String prefix);

class GeneratedRepositoryProvider extends SingleChildStatelessWidget {
  GeneratedRepositoryProvider({
    @required this.dio,
    this.prefixFn = _defaultPrefixFn,
  })  : assert(dio != null),
        assert(prefixFn != null);

  final Dio dio;
  final PrefixFn prefixFn;

  @override
  Widget buildWithChild(BuildContext context, Widget child) {
    return MultiProvider(
      providers: [
        Provider<OwnerRepository>(
          create: (_) => OwnerRepository(dio: dio, url: prefixFn('owner')),
        ),
        Provider<PetRepository>(
          create: (_) => PetRepository(dio: dio, url: prefixFn('pet')),
        ),
      ],
      child: child,
    );
  }
}

String _defaultPrefixFn(String url) => url;
