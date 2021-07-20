// GENERATED CODE - DO NOT MODIFY BY HAND
import 'dart:convert';

import 'package:json_annotation/json_annotation.dart';

import './pet.dart';

part 'owner.g.dart';

@JsonSerializable()
class Owner {
  Owner();

  @JsonKey(name: 'id')
  int? id;

  @JsonKey(name: 'name')
  String? name;

  @JsonKey(name: 'age')
  int? age;

  OwnerEdges? edges;

  @override
  int get hashCode => id.hashCode;

  @override
  bool operator ==(Object other) => other is Owner && id == other.id;

  factory Owner.fromJson(Map<String, dynamic> json) => _$OwnerFromJson(json);

  Map<String, dynamic> toJson() => _$OwnerToJson(this);

  String toString() => jsonEncode(toJson());
}

@JsonSerializable()
class OwnerEdges {
  OwnerEdges();

  List<Pet>? pets;

  factory OwnerEdges.fromJson(Map<String, dynamic> json) => _$OwnerEdgesFromJson(json);

  Map<String, dynamic> toJson() => _$OwnerEdgesToJson(this);
}
