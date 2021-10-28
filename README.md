    -ip string  
      	Search for a specific IP within all networks (default "none")  
    -net string  
      	Search for a specific CIDR block (default "none")  
    -reg string  
      	AWS region to return data for (default "none")  
    -ser string  
      	AWS Service to generate addresses for (ec2, aws, cf, r53, or all (default "none")  
    -v4  
      	Return only IPv4 values  
    -v6    
      	Return only IPv6 values    
  
  
    ./awsip -ser s3 -reg us-east-1  
    IP Prefix                                    Region              Service  
    =========================================================================    
    54.231.0.0/17                                us-east-1           S3  
    52.92.16.0/20                                us-east-1           S3  
    52.216.0.0/15                                us-east-1           S3  
    2600:1fa0:8000::/40                          us-east-1           S3  
    2600:1ffa:8000::/40                          us-east-1           S3  
    2600:1ff8:8000::/40                          us-east-1           S3  
    2600:1ff9:8000::/40                          us-east-1           S3  
  
    ./awsip -ser s3 -reg us-east-1 -v4  
    IP Prefix                                    Region              Service  
    =========================================================================    
    54.231.0.0/17                                us-east-1           S3  
    52.92.16.0/20                                us-east-1           S3  
    52.216.0.0/15                                us-east-1           S3  
  
    ./awsip -ser s3 -reg us-east-1 -v6  
    IP Prefix                                    Region              Service  
    =========================================================================  
    2600:1fa0:8000::/40                          us-east-1           S3  
    2600:1ffa:8000::/40                          us-east-1           S3  
    2600:1ff8:8000::/40                          us-east-1           S3  
    2600:1ff9:8000::/40                          us-east-1           S3  
