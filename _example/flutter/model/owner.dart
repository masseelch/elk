// GENERATED CODE - DO NOT MODIFY BY HAND
import 'dart:convert';

import 'package:json_annotation/json_annotation.dart';

import '../date_utc_converter.dart';

import '../color.dart';

import '../model/pet.dart';

part 'owner.g.dart';

@JsonSerializable()
@DateUtcConverter()
class Owner {
  Owner();

  int? id;
  String? name;

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

  factory OwnerEdges.fromJson(Map<String, dynamic> json) =>
      _$OwnerEdgesFromJson(json);
  Map<String, dynamic> toJson() => _$OwnerEdgesToJson(this);
}
