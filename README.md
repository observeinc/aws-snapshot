# aws-snapshot

This repository provides a library for scraping resource information from AWS.
It takes as an input a set of actions to execute, and outputs a list of events.

By design, we reuse the list of allowed actions, as defined in an IAM policy,
to drive the decision of what to scrape. This reduces the surface area for
misconfiguration significantly.

## Running aws-snapshot

A utility is provided that can be run as a standalone for debug purposes. You can build it through:

```
make
```

And execute the appropriate build for your local environment under `bin/`:

```
→ ./aws-snapshot -h
Usage of ./aws-snapshot:
  -buffer-size int
        Length of buffer for records (default 100)
  -manifest-file string
        Manifest filename
  -max-concurrent-requests int
        Maximum concurrent requests (default 10)
```

The utility will use your AWS profile credentials to scrape data. You can constrain collection by providing a manifest, e.g:

```
→ cat manifest.yaml
{
    "include": [ "rds:*", "s3:*" ],
    "exclude": [ "rds:DescribeDBInstances" ]
}
→ ./aws-snapshot -mainfest-file manifest.yaml
```
