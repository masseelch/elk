// GENERATED CODE - DO NOT MODIFY BY HAND
import 'package:dio/dio.dart';
import 'package:flutter/widgets.dart';
import 'package:json_annotation/json_annotation.dart';
import 'package:provider/provider.dart';

import '../date_utc_converter.dart';

import '../color.dart';

import '../model/pet.dart';
import '../model/owner.dart';
import '../client/owner.dart';

part 'pet.g.dart';

const petUrl = 'pets';

class PetClient {
  PetClient({required this.dio});

  final Dio dio;

  Future<Pet> find(int id) async {
    final r = await dio.get('/$petUrl/$id');
    return Pet.fromJson(r.data);
  }

  Future<List<Pet>> list({
    int? page,
    int? itemsPerPage,
    String? name,
    int? age,
    Color? color,
    List<String>? stringList,
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

    if (color != null) {
      params['color'] = color;
    }

    if (stringList != null) {
      params['string_list'] = stringList;
    }

    final r = await dio.get('/$petUrl');

    if (r.data == null) {
      return [];
    }

    return (r.data as List).map((i) => Pet.fromJson(i)).toList();
  }

  Future<Pet> create(PetCreateRequest req) async {
    final r = await dio.post('/$petUrl', data: req.toJson());
    return (Pet.fromJson(r.data));
  }

  Future<Pet> update(PetUpdateRequest req) async {
    final r = await dio.patch('/$petUrl/${req.id}', data: req.toJson());
    return (Pet.fromJson(r.data));
  }

  Future delete(int id) => dio.delete('/$petUrl/$id');

  Future<Owner> owner(Pet e) async {
    final r = await dio.get('/$petUrl/${e.id}/$ownerUrl');
    return (Owner.fromJson(r.data));
  }

  static PetClient of(BuildContext context) =>
      Provider.of<PetClient>(context, listen: false);
}

@JsonSerializable(createFactory: false)
@DateUtcConverter()
class PetCreateRequest {
  PetCreateRequest({
    this.name,
    this.age,
    this.color,
    this.stringList,
    this.owner,
  });

  PetCreateRequest.fromPet(Pet e)
      : name = e.name,
        age = e.age,
        color = e.color,
        stringList = e.stringList,
        owner = e.edges?.owner?.id;

  String? name;
  int? age;
  Color? color;
  List<String>? stringList;
  int? owner;

  Map<String, dynamic> toJson() => _$PetCreateRequestToJson(this);
}

@JsonSerializable(createFactory: false)
@DateUtcConverter()
class PetUpdateRequest {
  PetUpdateRequest({
    this.id,
    this.name,
    this.age,
    this.color,
    this.stringList,
    this.owner,
  });

  PetUpdateRequest.fromPet(Pet e)
      : id = e.id,
        name = e.name,
        age = e.age,
        color = e.color,
        stringList = e.stringList,
        owner = e.edges?.owner?.id;

  int id;
  String? name;
  int? age;
  Color? color;
  List<String>? stringList;
  int? owner;

  Map<String, dynamic> toJson() => _$PetUpdateRequestToJson(this);
}
