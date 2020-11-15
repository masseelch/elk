// GENERATED CODE - DO NOT MODIFY BY HAND
import 'package:json_annotation/json_annotation.dart';

part 'skip_generation_model.g.dart';

@JsonSerializable()
class SkipGenerationModel {
  SkipGenerationModel();

  int id;
  String name;

  SkipGenerationModelEdges edges;

  factory SkipGenerationModel.fromJson(Map<String, dynamic> json) =>
      _$SkipGenerationModelFromJson(json);
  Map<String, dynamic> toJson() => _$SkipGenerationModelToJson(this);
}

@JsonSerializable()
class SkipGenerationModelEdges {
  SkipGenerationModelEdges();

  factory SkipGenerationModelEdges.fromJson(Map<String, dynamic> json) =>
      _$SkipGenerationModelEdgesFromJson(json);
  Map<String, dynamic> toJson() => _$SkipGenerationModelEdgesToJson(this);
}
