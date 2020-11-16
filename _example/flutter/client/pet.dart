// GENERATED CODE - DO NOT MODIFY BY HAND
import 'package:dio/dio.dart';
import 'package:flutter/widgets.dart';
import 'package:provider/provider.dart';

import '../model/pet.dart';
import '../model/owner.dart';
import '../client/owner.dart';

const petUrl = 'pets';

class PetClient {
  PetClient({@required this.dio}) : assert(dio != null);

  final Dio dio;

  Future<Pet> find(int id) async {
    final r = await dio.get('/$petUrl/$id');
    return Pet.fromJson(r.data);
  }

  Future<List<Pet>> list({
    int page,
    int itemsPerPage,
    String name,
    int age,
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

    if (age != null) {
      params['age'] = age;
    }

    final r = await dio.get('/$petUrl');

    if (r.data == null) {
      return [];
    }

    return (r.data as List).map((i) => Pet.fromJson(i)).toList();
  }

  Future<Pet> create(Pet e) async {
    final r = await dio.post('/$petUrl', data: e.toJson());
    return (Pet.fromJson(r.data));
  }

  Future<Pet> update(Pet e) async {
    final r = await dio.patch('/$petUrl', data: e.toJson());
    return (Pet.fromJson(r.data));
  }

  Future<Owner> owner(Pet e) async {
    final r = await dio.get('/$petUrl/${e.id}/$ownerUrl');
    return (Owner.fromJson(r.data));
  }

  static PetClient of(BuildContext context) =>
      Provider.of<PetClient>(context, listen: false);
}
