// GENERATED CODE - DO NOT MODIFY BY HAND
import 'package:dio/dio.dart';
import 'package:flutter/widgets.dart';
import 'package:json_annotation/json_annotation.dart';
import 'package:provider/provider.dart';

import '../date_utc_converter.dart';

import '../color.dart';

import '../model/owner.dart';
import '../model/pet.dart';
import '../client/pet.dart';

part 'owner.g.dart';

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

  Future<Owner> create(OwnerCreateRequest req) async {
    final r = await dio.post('/$ownerUrl', data: req.toJson());
    return (Owner.fromJson(r.data));
  }

  Future<Owner> update(OwnerUpdateRequest req) async {
    final r = await dio.patch('/$ownerUrl/${req.id}', data: req.toJson());
    return (Owner.fromJson(r.data));
  }

  Future<List<Pet>> pets(Owner e) async {
    final r = await dio.get('/$ownerUrl/${e.id}/$petUrl');
    return (r.data as List).map((i) => Pet.fromJson(i)).toList();
  }

  static OwnerClient of(BuildContext context) =>
      Provider.of<OwnerClient>(context, listen: false);
}

@JsonSerializable(createFactory: false)
@DateUtcConverter()
class OwnerCreateRequest {
  OwnerCreateRequest({
    this.name,
    this.pets,
  });

  OwnerCreateRequest.fromOwner(Owner e)
      : name = e.name,
        pets = e.edges?.pets?.map((e) => e.id)?.toList();

  String name;
  List<int> pets;

  Map<String, dynamic> toJson() => _$OwnerCreateRequestToJson(this);
}

@JsonSerializable(createFactory: false)
@DateUtcConverter()
class OwnerUpdateRequest {
  OwnerUpdateRequest({
    this.id,
    this.name,
    this.pets,
  });

  OwnerUpdateRequest.fromOwner(Owner e)
      : id = e.id,
        name = e.name,
        pets = e.edges?.pets?.map((e) => e.id)?.toList();

  int id;
  String name;
  List<int> pets;

  Map<String, dynamic> toJson() => _$OwnerUpdateRequestToJson(this);
}
