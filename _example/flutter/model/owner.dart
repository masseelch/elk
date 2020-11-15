// GENERATED CODE - DO NOT MODIFY BY HAND
import 'package:json_annotation/json_annotation.dart';

import '../model/pet.dart';

part 'owner.g.dart';

@JsonSerializable()
class Owner {
  Owner();

  int id;
  String name;

  OwnerEdges edges;

  factory Owner.fromJson(Map<String, dynamic> json) => _$OwnerFromJson(json);
  Map<String, dynamic> toJson() => _$OwnerToJson(this);
}

@JsonSerializable()
class OwnerEdges {
  OwnerEdges();

  List<Pet> pets;

  factory OwnerEdges.fromJson(Map<String, dynamic> json) =>
      _$OwnerEdgesFromJson(json);
  Map<String, dynamic> toJson() => _$OwnerEdgesToJson(this);
}
