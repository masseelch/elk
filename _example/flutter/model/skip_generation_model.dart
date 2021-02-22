// GENERATED CODE - DO NOT MODIFY BY HAND
import 'dart:convert';

import 'package:json_annotation/json_annotation.dart';

import '../date_utc_converter.dart';

import '../color.dart';

part 'skip_generation_model.g.dart';

@JsonSerializable()
@DateUtcConverter()
class SkipGenerationModel {
  SkipGenerationModel();

  int? id;
  String? name;

  SkipGenerationModelEdges? edges;

  @override
  int get hashCode => id.hashCode;

  @override
  bool operator ==(Object other) =>
      other is SkipGenerationModel && id == other.id;

  factory SkipGenerationModel.fromJson(Map<String, dynamic> json) =>
      _$SkipGenerationModelFromJson(json);
  Map<String, dynamic> toJson() => _$SkipGenerationModelToJson(this);

  String toString() => jsonEncode(toJson());
}

@JsonSerializable()
class SkipGenerationModelEdges {
  SkipGenerationModelEdges();

  factory SkipGenerationModelEdges.fromJson(Map<String, dynamic> json) =>
      _$SkipGenerationModelEdgesFromJson(json);
  Map<String, dynamic> toJson() => _$SkipGenerationModelEdgesToJson(this);
}
