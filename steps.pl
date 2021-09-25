#!/usr/bin/perl -w
# Render the steps as HTML templates
use v5.30;
use Data::Dumper;
use standard;

sub renderPage {
    my ($n, $type_name, $value) = @_;
    say qq(<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<title>Grammatisch</title>
</head>
<body>
<form action="/step@{[sprintf "%02d", $n + 1]}" method="post">
<div>
<label for="$type_name">$type_name</label>);
    if (ref($value) eq 'ARRAY') { # is option
        say qq[<select id="$type_name" name="$type_name">];
        for ($value->@*) {
            say qq[<option value="$_">$_</option>];
        }
        say '</select>';
    } else {
        say qq[<textarea id="$type_name" name="$type_name"> {{ .$type_name }} </textarea>];
    }
    say '</div>
<input type="submit"/>
</form>
</body>
</html>';
}
