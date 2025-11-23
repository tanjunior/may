// GENERATED FROM Go constants in backend/errors.go
// Keep in sync with backend; used by frontend for error-code checks and messages.

export const CodeInternalError = "INTERNAL_ERROR";
export const CodeProductNotFound = "PRODUCT_NOT_FOUND";
export const CodeProductsNotFound = "PRODUCTS_NOT_FOUND";
export const CodeInvalidRequest = "INVALID_REQUEST";
export const CodeInvalidID = "INVALID_ID";
export const CodePerPageTooLarge = "PER_PAGE_TOO_LARGE";

export const ErrorMessages: Record<string, string> = {
  [CodeInternalError]: "internal server error",
  [CodeProductNotFound]: "product not found",
  [CodeProductsNotFound]: "products not found",
  [CodeInvalidRequest]: "invalid request",
  [CodeInvalidID]: "invalid product id",
  [CodePerPageTooLarge]: "per_page exceeds maximum allowed",
};

export default {
  CodeInternalError,
  CodeProductNotFound,
  CodeProductsNotFound,
  CodeInvalidRequest,
  CodeInvalidID,
  CodePerPageTooLarge,
  ErrorMessages,
};
