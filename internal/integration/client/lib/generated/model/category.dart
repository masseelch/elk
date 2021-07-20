// GENERATED CODE - DO NOT MODIFY BY HAND
import 'dart:convert';

import 'package:json_annotation/json_annotation.dart';

import './pet.dart';

part 'category.g.dart';

@JsonSerializable()
class Category {
  Category();

  @JsonKey(name: 'id')
  int? id;

  @JsonKey(name: 'name')
  String? name;

  CategoryEdges? edges;

  @override
  int get hashCode => id.hashCode;

  @override
  bool operator ==(Object other) => other is Category && id == other.id;

  factory Category.fromJson(Map<String, dynamic> json) => _$CategoryFromJson(json);

  Map<String, dynamic> toJson() => _$CategoryToJson(this);

  String toString() => jsonEncode(toJson());
}

@JsonSerializable()
class CategoryEdges {
  CategoryEdges();

  List<Pet>? pets;

  factory CategoryEdges.fromJson(Map<String, dynamic> json) => _$CategoryEdgesFromJson(json);

  Map<String, dynamic> toJson() => _$CategoryEdgesToJson(this);
}
