// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'pet.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Pet _$PetFromJson(Map<String, dynamic> json) {
  return Pet()
    ..id = json['id'] as int?
    ..name = json['name'] as String?
    ..age = json['age'] as int?
    ..edges = json['edges'] == null
        ? null
        : PetEdges.fromJson(json['edges'] as Map<String, dynamic>);
}

Map<String, dynamic> _$PetToJson(Pet instance) => <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'age': instance.age,
      'edges': instance.edges,
    };

PetEdges _$PetEdgesFromJson(Map<String, dynamic> json) {
  return PetEdges()
    ..category = (json['category'] as List<dynamic>?)
        ?.map((e) => Category.fromJson(e as Map<String, dynamic>))
        .toList()
    ..owner = json['owner'] == null
        ? null
        : Owner.fromJson(json['owner'] as Map<String, dynamic>)
    ..friends = (json['friends'] as List<dynamic>?)
        ?.map((e) => Pet.fromJson(e as Map<String, dynamic>))
        .toList();
}

Map<String, dynamic> _$PetEdgesToJson(PetEdges instance) => <String, dynamic>{
      'category': instance.category,
      'owner': instance.owner,
      'friends': instance.friends,
    };
