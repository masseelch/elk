{{ define "dateUtcConverter" }}
    import 'package:json_annotation/json_annotation.dart';

    class DateUtcConverter implements JsonConverter<DateTime?, String?> {
        const DateUtcConverter();

        @override
        DateTime? fromJson(String? json) => json == null ? null : DateTime.parse(json);

        @override
        String? toJson(DateTime? object) => object?.toUtc().toIso8601String();
    }
{{ end }}