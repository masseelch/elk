// GENERATED CODE - DO NOT MODIFY BY HAND
import 'dart:convert';

import 'package:json_annotation/json_annotation.dart';

import '../date_utc_converter.dart';

import '../color.dart';

import '../model/owner.dart';

part 'pet.g.dart';

@JsonSerializable()
@DateUtcConverter()
class Pet {
  Pet();

  int id;
  String name;
  int age;
  @ColorConverter()
  Color color;

  PetEdges edges;

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

  Owner owner;

  factory PetEdges.fromJson(Map<String, dynamic> json) =>
      _$PetEdgesFromJson(json);
  Map<String, dynamic> toJson() => _$PetEdgesToJson(this);
}
