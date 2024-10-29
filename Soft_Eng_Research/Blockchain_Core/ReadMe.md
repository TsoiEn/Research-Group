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

### Student Login

#### Create Mock Login Section
- Implement student login requiring only `studentID` and `password`.
- Validate login credentials against mock data.

### Additional Notes

#### Student Profile Information
- Include `firstName`, `middleName`, `lastName`, `Age`, and `email` in the student profile but not in the login credentials.

