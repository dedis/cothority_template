[![Build Status](https://travis-ci.org/dedis/cothority_template.svg?branch=master)](https://travis-ci.org/dedis/cothority_template)

Navigation: [DEDIS](https://github.com/dedis/doc/tree/master/README.md) ::
Cothority Template


# Template for a new cothority protocol/service/app

Building on the ONet-library available at
https://github.com/dedis/onet, this
repo holds templates to build the different parts necessary for a cothority
addition:

* [protocol](protocol) - define an ephemeral, distributed, decentralized protocol
* [service](service) - create a long-term service that can spawn any number of protocols
* [app](app) - write an app that will interact with, or spawn, a cothority
* [simulation](simulation) - how to create a simulation of a protocol or service

This repo is geared towards PhD-students who want to add a new functionality to
the cothority by creating their own protocols, services, simulations or apps.

## Starting

Just for testing you can `go get github.com/dedis/cothority_template`. For setting
up a new protocol/service/simulation, we propose that you create a new personal
repository in your account and then copy over the necessary files. Then you
will need to replace all the `github.com/dedis/cothority_tempate` references
with the path of your repository, e.g. `github.com/foo/super_protocol` if your
account is `foo` and the repository is `super_protocol`.

If you happen to do a semester project for DEDIS, please ask your responsible to
set up a `github.com/dedis/student_yy_name` repository for you.

The Perl pie to the rescue (or `sed -i` if you prefer...):

```bash
find . -name "*go" | xargs \
perl -pi -e "s:github.com/dedis/cothority_template:github.com/foo/super_protocol:"
```

## Documentation

You find more documentation on how to use the template on the wiki:
[Cothority Template](https://github.com/dedis/cothority_template/wiki)

More documentation and examples can be found at:
- To run and use a conode, have a look at
	[Cothority Node](https://github.com/dedis/cothority)
	with examples of protocols, services and apps
- To participate as a core-developer, go to
	[Cothority Network Library](https://github.com/dedis/onet)

## License

All repositories for the cothority are double-licensed under a
GNU/AGPL 3.0 and a commercial license. If you want to have more information,
contact us at dedis@epfl.ch.

## Contribution

If you want to contribute to Cothority-ONet, please have a look at
[CONTRIBUTION](https://github.com/dedis/cothority/blob/master/CONTRIBUTION) for
licensing details. Once you are OK with those, you can have a look at our
coding-guidelines in
[Coding](https://github.com/dedis/Coding). In short, we use the github-issues
to communicate and pull-requests to do code-review. Travis makes sure that
everything goes smoothly. And we'd like to have good code-coverage.

# Contact

You can contact us at https://groups.google.com/forum/#!forum/cothority
