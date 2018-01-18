this example assumes you have run the crunchy-containers
example/docker/primary-replica example.  That example will
create a primary container and replica container that form
a postgres cluster. 

This watch example looks for a primary container, if not found, 
it will trigger a failover onto the replica container.

To test this example, delete the master pod and examine
the log of the watch container to see if perform the failover logic.
