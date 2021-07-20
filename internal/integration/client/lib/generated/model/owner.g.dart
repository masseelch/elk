// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'owner.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Owner _$OwnerFromJson(Map<String, dynamic> json) {
  return Owner()
    ..id = json['id'] as int?
    ..name = json['name'] as String?
    ..age = json['age'] as int?
    ..edges = json['edges'] == null
        ? null
        : OwnerEdges.fromJson(json['edges'] as Map<String, dynamic>);
}

Map<String, dynamic> _$OwnerToJson(Owner instance) => <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'age': instance.age,
      'edges': instance.edges,
    };

OwnerEdges _$OwnerEdgesFromJson(Map<String, dynamic> json) {
  return OwnerEdges()
    ..pets = (json['pets'] as List<dynamic>?)
        ?.map((e) => Pet.fromJson(e as Map<String, dynamic>))
        .toList();
}

Map<String, dynamic> _$OwnerEdgesToJson(OwnerEdges instance) =>
    <String, dynamic>{
      'pets': instance.pets,
    };
