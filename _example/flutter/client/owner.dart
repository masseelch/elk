// GENERATED CODE - DO NOT MODIFY BY HAND
import 'package:dio/dio.dart';
import 'package:flutter/widgets.dart';
import 'package:provider/provider.dart';

import '../model/owner.dart';
import '../model/pet.dart';
import '../client/pet.dart';

const ownerUrl = 'owners';

class OwnerClient {
  OwnerClient({@required this.dio}) : assert(dio != null);

  final Dio dio;

  Future<Owner> find(int id) async {
    final r = await dio.get('/$ownerUrl/$id');
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

    final r = await dio.get('/$ownerUrl');

    if (r.data == null) {
      return [];
    }

    return (r.data as List).map((i) => Owner.fromJson(i)).toList();
  }

  Future<Owner> create(Owner e) async {
    final r = await dio.post('/$ownerUrl', data: e.toJson());
    return (Owner.fromJson(r.data));
  }

  Future<Owner> update(Owner e) async {
    final r = await dio.patch('/$ownerUrl', data: e.toJson());
    return (Owner.fromJson(r.data));
  }

  Future<List<Pet>> pets(Owner e) async {
    final r = await dio.get('/$ownerUrl/${e.id}/$petUrl');
    return (r.data as List).map((i) => Pet.fromJson(i)).toList();
  }

  static OwnerClient of(BuildContext context) =>
      Provider.of<OwnerClient>(context, listen: false);
}
