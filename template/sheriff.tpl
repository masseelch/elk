{{ define "model/stringer" }}
    // Get rid of the fmt.Stringer implementation since it breaks liip/sheriff.
    // These lines have to be here since template/text does skip empty templates.
{{ end }}