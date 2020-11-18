// GENERATED CODE - DO NOT MODIFY BY HAND
import 'package:json_annotation/json_annotation.dart';

import '../color.dart';

import '../model/owner.dart';

part 'pet.g.dart';

@JsonSerializable()
class Pet {
  Pet();

  int id;
  String name;
  int age;
  dynamic color;

  PetEdges edges;

  factory Pet.fromJson(Map<String, dynamic> json) => _$PetFromJson(json);
  Map<String, dynamic> toJson() => _$PetToJson(this);
}

@JsonSerializable()
class PetEdges {
  PetEdges();

  Owner owner;

  factory PetEdges.fromJson(Map<String, dynamic> json) =>
      _$PetEdgesFromJson(json);
  Map<String, dynamic> toJson() => _$PetEdgesToJson(this);
}
