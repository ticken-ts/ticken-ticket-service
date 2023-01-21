// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

struct Ticket {
    string  section;
    bool    scanned;
}

struct TicketDTO {
    uint256 id;
    string  section;
    bool    scanned;
}

contract TickenEvent is ERC721, Pausable, Ownable {

    event TicketCreated(
        address ownerAddress,
        uint256 indexed tokenID,
        string section
    );

    using Counters for Counters.Counter;

    Counters.Counter private _tokenIdCounter;

    // Mapping from token ID to ticket
    mapping(uint256 => Ticket) private _tickets;

    // Mapping from owner to list of owned token IDs
    mapping(address => uint256[]) private _ownedTokens;

    constructor() ERC721("TickenEvent", "TE") {}

    function pause() public onlyOwner {
        _pause();
    }

    function unpause() public onlyOwner {
        _unpause();
    }

    function safeMint(address to, string memory section) public whenNotPaused onlyOwner {
        uint256 tokenId = _tokenIdCounter.current();

        _tokenIdCounter.increment();

        _safeMint(to, tokenId);

        _tickets[tokenId] = Ticket(section, false);

        _ownedTokens[to].push(tokenId);

        emit TicketCreated(to, tokenId, section);
    }

    // Get all tickets from an owner address
    function getTicketsByOwner(address owner) public view returns (TicketDTO[] memory) {
        uint256[] memory tokenIds = _ownedTokens[owner];
        TicketDTO[] memory tickets = new TicketDTO[](tokenIds.length);
        for (uint256 i = 0; i < tokenIds.length; i++) {
            uint256 tokenId = tokenIds[i];
            Ticket memory ticket = _tickets[tokenId];
            tickets[i] = TicketDTO(tokenId, ticket.section, ticket.scanned);
        }
        return tickets;
    }

    function scanTicket(uint256 tokenId) public whenNotPaused onlyOwner {
        require(_exists(tokenId), "ERC721: operator query for nonexistent token");
        _tickets[tokenId].scanned = true;
    }

    function getTicket(uint256 tokenId) public view returns (Ticket memory) {
        require(_exists(tokenId), "ERC721: operator query for nonexistent token");
        return _tickets[tokenId];
    }

    function _beforeTokenTransfer(address from, address to, uint256 tokenId, uint256 batchSize)
    internal
    whenNotPaused
    override
    {
        super._beforeTokenTransfer(from, to, tokenId, batchSize);
    }
}
