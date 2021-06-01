#!/usr/bin/perl -w
use v5.30;
use Data::Dumper;
use Tie::IxHash;
use List::Util qw(pairs);
no warnings qw(experimental::smartmatch);

tie my %interface_srcs, 'Tie::IxHash';
tie my %struct_srcs, 'Tie::IxHash';
tie my %func_srcs, 'Tie::IxHash';
tie my %struct_members_refv, 'Tie::IxHash';
my @function_srcs;

sub unmarshal {
    my ($name) = @_; # struct_name
    my $ret = "func (x *$name) UnmarshalJSON(data []byte) error {\n";
    $ret .= tmpStruct($name) . "\n";
    $ret .= "var y tmp
json.Unmarshal(data, &y)\n";
    for ($struct_members_refv{$name}->@*) {
        my ($name, $type) = @$_;
        $ret .= ($interface_srcs{$type} ?
            unmarshalInterfaceAssign($name, $type) :
            structAssign($name)) . "\n";
    }
    $ret .= "return nil
}";
    return $ret;
}

sub possibleStructs {
    my $ret = "possibleStructs := map[string][]interface{}\{\n";
    for (keys %interface_srcs) {
        my @impls = implementations($_);
        $ret .= qq["$_": []$_\{\n];
        for my $impl (@impls) {
            $ret .= "&$impl\{},\n";
        }

        $ret .= "},\n";
    }
    $ret .= '}';
    return $ret;
}

sub implementations {
    my ($iface_name) = @_;
    my $suffix = substr $iface_name, 1;
    return grep /$suffix/, keys %struct_srcs;
}

sub tmpStruct {
    my ($struct_name) = @_;
    my $ret = "type tmp struct {\n";
    for ($struct_members_refv{$struct_name}->@*) {
        my ($name, $type) = @$_;
        $ret .= "$name " .
            ($interface_srcs{$type} ? 'json.RawMessage' : $type) . "\n";
    }
    $ret .= '}';
    return $ret;
}

sub unmarshalInterfaceAssign {
    my ($name, $type) = @_;
    my $instance_name = "\l${type}Instance";
    my $ret = "var $instance_name map[string]interface{}\n";
    $ret .= "json.Unmarshal(y.$name, &$instance_name)\n";
    for (implementations($type)) {
        $ret .= qq<if $instance_name\["Type"] == "$_" {
var z $_
json.Unmarshal(y.$name, &z)
x.$name = &z
}\n>;
    }
    chomp $ret;
    return $ret;
}

sub structAssign {
    my ($name) = @_;
    return "x.$name = y.$name";
}

sub addMember {
    my ($src, $member) = @_;
    substr $src, index($src, "\n"), 0, "\n$member";
    return $src;
}

sub isValid {
    my ($struct) = @_;
    return qq[func (x *$struct) isValid() bool {
return x.Type == "$struct";
}];
}

$/ = '';
while (<DATA>) {
    chomp;
    if (/^type (\w+) interface \{$/m) {
        $interface_srcs{$1} = $_;
    } elsif (/^type (\w+) struct \{$/m) {
        my $struct_name = $1;
        $struct_srcs{$struct_name} = $_;
        $struct_members_refv{$struct_name} = [pairs split ' ', s/^.*\n((?:.*\n)*)\s*\}$/$1/r];
    } elsif (/^func/) {
        push @function_srcs, $_;
    } else {
        die "Paragraph not recognized: $_";
    }
}

say qq[package main
import (
"encoding/json"
)];
while (my ($_k, $v) = each %interface_srcs) {
    say $v;
}

my @struct_srcs_keys = keys %struct_srcs;
for my $k (@struct_srcs_keys) {
    my $v = $struct_srcs{$k};
    say addMember $v, 'Type string';
    say unmarshal $k;
}

for (@function_srcs) {
    say $_;
}


__DATA__
type Op interface {
	ToRegex()
}

type createNodesOp struct {
	Count int
    CountParam Op
}

func (*createNodesOp) ToRegex() {
}

type createNamespacesOp struct {
	Prefix string
}

func (*createNamespacesOp) ToRegex() {
}

type createPodsOp struct {
	CollectMetrics bool
    Namespace *string
}

func (*createPodsOp) ToRegex() {
}

type Foo interface {
}

type barFoo struct {
}

type bazFoo struct {
}
