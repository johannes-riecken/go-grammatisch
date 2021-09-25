#!/usr/bin/perl -w
use v5.20;
use JSON::PP;
use re 'eval';
use standard;
my $text;
open my $f_regex, '<', 'regex.txt';
{ local $/; $text = <$f_regex>; }
close $f_regex;
my $re = qr/$text/x;
open my $f_input, '<', 'input.txt';
{ local $/; $_ = <$f_input>; }
close $f_input;
chomp;
/$re/;
say encode_json($^R->[1]);
