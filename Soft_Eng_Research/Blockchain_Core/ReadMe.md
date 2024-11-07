# Project Documentation

## Overview

This repository provides a checklist and documentation for the following knowledge areas and algorithms related to blockchain technology.

## Knowledge Needs

To effectively understand and work with blockchain technology, ensure you are familiar with the following topics:

1. **Blockchain Fundamentals**: Understand the core principles and architecture of blockchain technology.
2. **Data Structures**: Learn about the data structures commonly used in blockchain implementations.
3. **Cryptography**: Gain knowledge about cryptographic techniques used for securing transactions and data.
4. **Transaction Formats**: Familiarize yourself with the different formats of transactions in blockchain systems.
5. **Merkle Tree**: Study the construction and application of Merkle Trees in blockchain.
6. **Consensus Algorithm**: Learn about various consensus algorithms that help achieve agreement across distributed systems.

## Algorithms

This section outlines key algorithms relevant to blockchain technology:

1. **Hashing**: Learn about hashing algorithms used to ensure data integrity and create unique identifiers.
2. **Merkle Tree Construction**: Understand the process of constructing Merkle Trees for efficient and secure data verification.
3. **Block Linking**: Study how blocks are linked together to form a secure chain.
4. **Transaction Validation**: Learn how transactions are validated and verified within a blockchain network.




## StuCred Blockchain Project: Admin & Student Login

### Admin Functionality

#### Define Credential Types
- Update `credential.go` to include academic and non-academic credential distinctions.

#### Implement Access Controls in `chaincode.go`
- Allow admin (developer) to add, update, and delete academic credentials.
- Restrict admin access to non-academic credentials (students only).

#### Separate Methods for Each Credential Type in `chaincode.go`
- Add functions such as `AddAcademicCredential` and `UpdateAcademicCredential` for admin use.
- Implement `AddNonAcademicCredential` for students to manage their own credentials.

### Student Login (Finish)

#### Create Mock Login Section
- Implement student login requiring only `studentID` and `password`.
- Validate login credentials against mock data.

### Additional Notes

#### Student Profile Information
- Include `firstName`, `middleName`, `lastName`, `Age`, and `email` in the student profile but not in the login credentials.


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
│   │   │   ├── credential.go
│   │   │   ├── student.go
│   │   │   ├── admin.go         # Optional, for admin-specific features
│   │   ├── go.mod               # Module dependencies specific to chaincode
│   │   ├── go.sum               # Dependency checksum file
├── main.go                      # Used for testing but not part of the deployed chaincode
├── fabric-config/               # Network configuration files
│   ├── configtx.yaml            # Channel and orderer configurations
│   ├── crypto-config.yaml       # Crypto material generation definition
│   ├── docker-compose.yaml      # Docker orchestration for Fabric network
│   ├── genesis.block            # Genesis block for the orderer
│   ├── mychannel.tx             # Channel creation transaction file
├── crypto-config/               # Generated crypto material (certificates and keys)
│   ├── ordererOrganizations/
│   ├── peerOrganizations/
├── orderers/                    # Orderer-specific configurations
│   ├── orderer1/
│   │   ├── orderer.yaml         # Configuration specific to orderer1
│   │   ├── crypto/              # Crypto material for orderer1
│   ├── orderer2/
│   │   ├── orderer.yaml         # Configuration specific to orderer2
│   │   ├── crypto/              # Crypto material for orderer2
├── tools/                       # Optional tools (CouchDB, Explorer, etc.)
│   ├── couchdb/
│   ├── explorer/
│   │   ├── config.json
├── other_project_files/         # Directory for other project files, scripts, or configuration


