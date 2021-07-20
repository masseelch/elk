// GENERATED CODE - DO NOT MODIFY BY HAND
import 'dart:convert';

import 'package:json_annotation/json_annotation.dart';

import './category.dart';
import './owner.dart';

part 'pet.g.dart';

@JsonSerializable()
class Pet {
  Pet();

  @JsonKey(name: 'id')
  int? id;

  @JsonKey(name: 'name')
  String? name;

  @JsonKey(name: 'age')
  int? age;

  PetEdges? edges;

  @override
  int get hashCode => id.hashCode;

  @override
  bool operator ==(Object other) => other is Pet && id == other.id;

  factory Pet.fromJson(Map<String, dynamic> json) => _$PetFromJson(json);

  Map<String, dynamic> toJson() => _$PetToJson(this);

  String toString() => jsonEncode(toJson());
}

@JsonSerializable()
class PetEdges {
  PetEdges();

  List<Category>? category;
  Owner? owner;
  List<Pet>? friends;

  factory PetEdges.fromJson(Map<String, dynamic> json) => _$PetEdgesFromJson(json);

  Map<String, dynamic> toJson() => _$PetEdgesToJson(this);
}
