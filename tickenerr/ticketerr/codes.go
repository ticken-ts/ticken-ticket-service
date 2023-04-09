package ticketerr

const (
	FailedToGenerateTokenIDErrorCode = iota + 300
	FailedToMintTicketInPVTBC
	FailedToMintTicketInPUBBC
	FailedToRetrievePUBBCTicketsErrorCode
	TicketNotFoundErrorCode
	BuyResellErrorCode
	CreateResellErrorCode
	ResellCurrencyNotSupportedErrorCode
	TicketResellNotFoundErrorCode
	TransferTicketInPUBBCErrorCode
	ReadUserTicketsFromDatabaseErrorCode
)

func GetErrorMessages(code uint32) string {
	switch code {
	case FailedToGenerateTokenIDErrorCode:
		return "an error occurred while generating ticket associated token ID"
	case FailedToMintTicketInPVTBC:
		return "failed to mint ticket in private blockchain"
	case FailedToMintTicketInPUBBC:
		return "failed to mint ticket in public blockchain"
	case FailedToRetrievePUBBCTicketsErrorCode:
		return "failed to retrieve ticket from public blockchain"
	case TicketNotFoundErrorCode:
		return "ticket not found"
	case BuyResellErrorCode:
		return "failed to buy resell"
	case CreateResellErrorCode:
		return "failed to create ticket resell"
	case ResellCurrencyNotSupportedErrorCode:
		return "the currency used to create the resell is not currently supported"
	case TicketResellNotFoundErrorCode:
		return "ticket resell not found"
	case TransferTicketInPUBBCErrorCode:
		return "an error occurred while transferring the ticket in the public blockchain"
	case ReadUserTicketsFromDatabaseErrorCode:
		return "an error occurred while trying to read user tickets from database"
	default:
		return "an error has occurred"
	}
}
