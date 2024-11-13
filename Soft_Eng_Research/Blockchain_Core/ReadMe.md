Name: Hank B Davis, Age: 20, ID: 202533282, Email: dhb3282@example.edu.ph, Password: zpKVM4cQ
Name: Charlie H Jones, Age: 25, ID: 202403450, Email: jch3450@example.edu.ph, Password: S2oS6gQP
Name: Bob I Rodriguez, Age: 20, ID: 202209675, Email: rbi9675@example.edu.ph, Password: WXwoXDA9
Name: Alice D Smith, Age: 25, ID: 202433194, Email: sad3194@example.edu.ph, Password: uOlgXCpt
Name: Diana F Garcia, Age: 18, ID: 202226488, Email: gdf6488@example.edu.ph, Password: TOTVufqI
Name: Alice J Garcia, Age: 18, ID: 202413171, Email: gaj3171@example.edu.ph, Password: 4G2mn3Fx
Name: Eve D Brown, Age: 19, ID: 202120988, Email: bed0988@example.edu.ph, Password: IeGiLref
Name: Bob C Garcia, Age: 20, ID: 202329393, Email: gbc9393@example.edu.ph, Password: ndIfT6TM
Name: Frank J Jones, Age: 24, ID: 202207626, Email: jfj7626@example.edu.ph, Password: v2hXvvv4
Name: Frank E Martinez, Age: 24, ID: 202203708, Email: mfe3708@example.edu.ph, Password: Sz6AkUMD


### current structure of blockchain
blockchain/
├── chaincode/
│   ├── src/
│   │   ├── chaincode.go         # Main chaincode logic
│   │   ├── model/
│   │   │   ├── block.go         # Model files as dependencies
│   │   │   ├── admin.go
│   │   │   ├── credential.go
│   │   │   ├── student.go
│   │   │   ├── utils.go         # Optional utility functions
│   │   ├── go.mod               # Module dependencies specific to chaincode
│   │   ├── go.sum
├── consensus/                  # Raft-related files
│   ├── raft.go                 # Main Raft logic
│   ├── raft_node.go            # Raft node management
│   ├── raft_helpers.go         # Helper functions for Raft
│   └── raft_config.go          # Raft configuration setup
├── api/
                  # Used for testing but 


### admin.go responsibilities

- Adding new students
- Adding academic credentials
- Managing operations overseen by an admin, such as:
    - Overseeing blockchain updates
    - User management

### block.go
- manages block creation
- hashing
- serialization

### credential.go and student.go 
- handle data models related to credentials and students, respectively.

### chaincode.go 
- is the main entry point where chaincode logic interacts with the blockchain.



### API Integration

To connect the frontend and backend, you should create a new directory for your API within the `blockchain/` directory. Here is a suggested structure:

```
blockchain/
├── chaincode/
│   ├── src/
│   │   ├── chaincode.go         # Main chaincode logic
│   │   ├── model/
│   │   │   ├── block.go         # Model files as dependencies
│   │   │   ├── admin.go
│   │   │   ├── credential.go
│   │   │   ├── student.go
│   │   │   ├── admin.go         # Optional, for admin-specific features
│   │   ├── go.mod               # Module dependencies specific to chaincode
│   │   ├── go.sum               # Dependency checksum file
├── api/                         # New directory for API
│   ├── main.go                  # Entry point for the API server
│   ├── handlers/
│   │   ├── studentHandler.go    # Handlers for student-related API endpoints
│   │   ├── adminHandler.go      # Handlers for admin-related API endpoints
│   ├── router.go                # Router configuration
│   ├── go.mod                   # Module dependencies for the API
│   ├── go.sum                   # Dependency checksum file
```

This structure keeps your API logic separate from your chaincode logic, making it easier to manage and maintain.