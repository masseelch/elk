// GENERATED CODE - DO NOT MODIFY BY HAND
import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';

import '../model/owner.dart';

class OwnerRepository {
  OwnerRepository(
    @required this.dio,
    @required this.url,
  )   : assert(dio != null),
        assert(url != null && url != '');

  final Dio dio;
  final String url;

  Future<Owner> find(int id) async {
    final r = await dio.get('/$url/$id');
    return Owner.fromJson(r.data);
  }

  Future<List<Owner>> list({
    int page,
    int itemsPerPage,
    String name,
  }) async {
    final params = const {};

    if (page != null) {
      params['page'] = page;
    }

    if (itemsPerPage != null) {
      params['itemsPerPage'] = itemsPerPage;
    }

    if (name != null) {
      params['name'] = name;
    }

    final r = await dio.get('/$url');

    if (r.data == null) {
      return [];
    }

    return (r.data as List).map((i) => Owner.fromJson(i)).toList();
  }

  Future<Owner> create(Owner e) async {
    final r = await dio.post('/$url', data: e.toJson());
    return (Owner.fromJson(r.data));
  }

  Future<Owner> update(Owner e) async {
    final r = await dio.patch('/$url', data: e.toJson());
    return (Owner.fromJson(r.data));
  }
}
