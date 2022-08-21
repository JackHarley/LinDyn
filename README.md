# LinDyn

LinDyn is a simple dynamic DNS updater for Linode hosted DNS records.

It uses https://www.my-ip.io to retrieve the current IPv4 address and then checks if the configured Linode DNS zone A records are set to that address. If a mismatch is detected, the program updates the records automatically.

Configuration
------------------------------
LinDyn requires three environment variables:
* `LINODE_PERSONAL_ACCESS_TOKEN` - Your Linode Personal Access Token, generate one here: https://cloud.linode.com/profile/tokens. The only permission required is Domain: Read/Write.
* `LINODE_DOMAIN` - Name of your Linode domain zone as shown in the first column on the Domains page: https://cloud.linode.com/domains
* `LINODE_A_RECORDS` - Comma separated list of A record hostnames to ensure have their target set to the executor's current IPv4 address. The empty string typically represents the domain itself, e.g. if your domain zone is for example.com, the empty string record is for example.com.

Note: New records will **not** be created by LinDyn. It will only update existing records.

You can set the environment variables in any manner you choose, or you can place them in a file called `.env` in the directory where you will execute the program and LinDyn will read them from said file.


Example Usage
------------------------------

```
LINODE_PERSONAL_ACCESS_TOKEN=abcd1234 \
    LINODE_DOMAIN=mydomain.com \
    LINODE_A_RECORDS=",*,home" \
    ./lindyn
```

Or if using a `.env` file:
```
./lindyn
```


Build Instructions
------------------------------
```
make
```

To build for MIPS64 (EdgeRouter 4):
```
make lindyn-mips64
```

License
------------------------------
This software is licensed under the MIT License. For more information see [LICENSE.md]().