% Lamina workspace knowledge base
% Modules and their dependency relationships

% module(Name) — a workspace module
module(axon).
module(axon_tool).
module(axon_loop).
module(axon_talk).
module(axon_lens).
module(axon_chat).
module(axon_memo).
module(axon_auth).
module(axon_gate).
module(axon_look).
module(axon_task).
module(axon_fact).
module(axon_mind).
module(axon_eval).

% kind(Module, Kind) — library or service
kind(axon, library).
kind(axon_tool, library).
kind(axon_loop, library).
kind(axon_talk, library).
kind(axon_lens, library).
kind(axon_fact, library).
kind(axon_mind, library).
kind(axon_chat, service).
kind(axon_memo, service).
kind(axon_auth, service).
kind(axon_gate, service).
kind(axon_look, service).
kind(axon_task, service).
kind(axon_eval, standalone).

% depends_on(Module, Dependency) — direct dependency
depends_on(axon_loop, axon_tool).
depends_on(axon_talk, axon_loop).
depends_on(axon_chat, axon).
depends_on(axon_chat, axon_loop).
depends_on(axon_chat, axon_tool).
depends_on(axon_memo, axon).
depends_on(axon_auth, axon).
depends_on(axon_gate, axon).
depends_on(axon_look, axon).
depends_on(axon_task, axon).
depends_on(axon_mind, axon).
depends_on(axon_lens, axon_tool).

% Transitive dependency
transitive_dep(X, Y) :- depends_on(X, Y).
transitive_dep(X, Y) :- depends_on(X, Z), transitive_dep(Z, Y).

% Modules affected by a change to a library
affected_by(Lib, Consumer) :- transitive_dep(Consumer, Lib).

% Modules with no dependencies
standalone(X) :- module(X), \+ depends_on(X, _).

% Services that depend on a given library
service_using(Lib, Service) :- kind(Service, service), transitive_dep(Service, Lib).
