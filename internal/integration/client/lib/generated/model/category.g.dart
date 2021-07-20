// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'category.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Category _$CategoryFromJson(Map<String, dynamic> json) {
  return Category()
    ..id = json['id'] as int?
    ..name = json['name'] as String?
    ..edges = json['edges'] == null
        ? null
        : CategoryEdges.fromJson(json['edges'] as Map<String, dynamic>);
}

Map<String, dynamic> _$CategoryToJson(Category instance) => <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'edges': instance.edges,
    };

CategoryEdges _$CategoryEdgesFromJson(Map<String, dynamic> json) {
  return CategoryEdges()
    ..pets = (json['pets'] as List<dynamic>?)
        ?.map((e) => Pet.fromJson(e as Map<String, dynamic>))
        .toList();
}

Map<String, dynamic> _$CategoryEdgesToJson(CategoryEdges instance) =>
    <String, dynamic>{
      'pets': instance.pets,
    };
