# QFields

## Overview

QFields were origonally written by John Zarate to make the writting of GALIoT Systems golang graphQL qeuries with objects in object in objects (etc) much easier. The code was added to github by another one of the API developers Samuel Archibald (IoTPanic). QFields are generated from the context `gqlgen` provides in a resolver which will contain the fields within that query, the arguments in the object, and the children objects to be queried.

## Example

An example on how to use QFields is provided at `https://github.com/IoTPanic/cluster-monitor-qfield-graphql-example`