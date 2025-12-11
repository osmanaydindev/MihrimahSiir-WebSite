import DOMPurify from 'dompurify'

/**
 * Composable for input sanitization and XSS prevention
 */
export function useSanitizer() {
  /**
   * Sanitize HTML content to prevent XSS attacks
   * @param dirty - The potentially unsafe HTML string
   * @param config - Optional DOMPurify configuration
   * @returns Sanitized HTML string
   */
  const sanitizeHTML = (
    dirty: string,
    config?: DOMPurify.Config
  ): string => {
    return DOMPurify.sanitize(dirty, config)
  }

  /**
   * Strip all HTML tags from input, leaving only text
   * @param input - The input string with potential HTML
   * @returns Plain text without HTML tags
   */
  const stripHTML = (input: string): string => {
    return DOMPurify.sanitize(input, { ALLOWED_TAGS: [] })
  }

  /**
   * Sanitize user input for display (allows basic formatting)
   * @param input - User input string
   * @returns Sanitized string with safe HTML tags
   */
  const sanitizeUserInput = (input: string): string => {
    return DOMPurify.sanitize(input, {
      ALLOWED_TAGS: ['b', 'i', 'em', 'strong', 'p', 'br'],
      ALLOWED_ATTR: [],
    })
  }

  /**
   * Sanitize poem or book content (allows more formatting)
   * @param content - Content string
   * @returns Sanitized content with safe formatting
   */
  const sanitizeContent = (content: string): string => {
    return DOMPurify.sanitize(content, {
      ALLOWED_TAGS: [
        'p', 'br', 'b', 'i', 'em', 'strong', 'u',
        'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
        'blockquote', 'ul', 'ol', 'li'
      ],
      ALLOWED_ATTR: [],
    })
  }

  /**
   * Check if string contains dangerous content
   * @param input - Input string to check
   * @returns true if dangerous content detected
   */
  const containsDangerousContent = (input: string): boolean => {
    const lowerInput = input.toLowerCase()

    // Check for script tags
    if (lowerInput.includes('<script') || lowerInput.includes('</script>')) {
      return true
    }

    // Check for javascript: protocol
    if (lowerInput.includes('javascript:')) {
      return true
    }

    // Check for data: protocol (can be used for XSS)
    if (lowerInput.includes('data:text/html')) {
      return true
    }

    // Check for event handlers
    const dangerousPatterns = [
      'onerror', 'onload', 'onclick', 'onmouseover',
      'onfocus', 'onblur', 'onchange', 'onsubmit',
      '<iframe', '<object', '<embed'
    ]

    return dangerousPatterns.some(pattern => lowerInput.includes(pattern))
  }

  /**
   * Validate and sanitize string input
   * @param input - Input string
   * @param maxLength - Maximum allowed length
   * @returns Sanitized and validated string
   */
  const validateAndSanitize = (
    input: string,
    maxLength: number = 1000
  ): { valid: boolean; sanitized: string; error?: string } => {
    if (!input || input.trim() === '') {
      return { valid: false, sanitized: '', error: 'Input is required' }
    }

    if (containsDangerousContent(input)) {
      return {
        valid: false,
        sanitized: '',
        error: 'Content contains potentially dangerous elements'
      }
    }

    const sanitized = stripHTML(input).trim()

    if (sanitized.length > maxLength) {
      return {
        valid: false,
        sanitized: sanitized.substring(0, maxLength),
        error: `Content must not exceed ${maxLength} characters`
      }
    }

    return { valid: true, sanitized }
  }

  /**
   * Sanitize search query
   * @param query - Search query string
   * @returns Sanitized search query
   */
  const sanitizeSearchQuery = (query: string): string => {
    // Remove HTML tags
    let sanitized = stripHTML(query)

    // Remove special characters that could cause issues
    sanitized = sanitized.replace(/[<>]/g, '')

    // Limit length
    if (sanitized.length > 200) {
      sanitized = sanitized.substring(0, 200)
    }

    return sanitized.trim()
  }

  /**
   * Escape special characters for safe insertion into HTML
   * @param text - Text to escape
   * @returns Escaped text
   */
  const escapeHTML = (text: string): string => {
    const div = document.createElement('div')
    div.textContent = text
    return div.innerHTML
  }

  return {
    sanitizeHTML,
    stripHTML,
    sanitizeUserInput,
    sanitizeContent,
    containsDangerousContent,
    validateAndSanitize,
    sanitizeSearchQuery,
    escapeHTML,
  }
}
