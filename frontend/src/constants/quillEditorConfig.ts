/**
 * QuillEditor toolbar configurations
 * Reusable configurations for different use cases
 */

/**
 * Basic toolbar for simple text formatting (titles, short content)
 * Includes: text formatting, colors, alignment
 */
export const QUILL_TOOLBAR_BASIC = [
  ['bold', 'italic', 'underline', 'strike'],
  [{ 'color': [] }, { 'background': [] }],
  [{ 'align': [] }],
  ['clean']
]

/**
 * Full toolbar for rich content editing
 * Includes: text formatting, headers, lists, colors, alignment, links
 */
export const QUILL_TOOLBAR_FULL = [
  ['bold', 'italic', 'underline', 'strike'],
  [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
  [{ 'list': 'ordered'}, { 'list': 'bullet' }],
  [{ 'color': [] }, { 'background': [] }],
  [{ 'align': [] }],
  ['link'],
  ['clean']
]

/**
 * Minimal toolbar for very simple formatting
 * Includes: basic text formatting only
 */
export const QUILL_TOOLBAR_MINIMAL = [
  ['bold', 'italic', 'underline', 'strike']
]

/**
 * Default QuillEditor props for consistent configuration
 */
export const QUILL_DEFAULT_PROPS = {
  theme: 'snow' as const,
  contentType: 'html' as const
}
