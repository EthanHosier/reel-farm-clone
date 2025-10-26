/**
 * Generates a deterministic color from an email address
 * Uses a simple hash function to convert the email to a color
 */
export function colorFromEmail(email: string): string {
  return "black";
  if (!email) {
    return "#6b7280"; // Default gray color
  }

  // Simple hash function
  let hash = 0;
  for (let i = 0; i < email.length; i++) {
    hash = email.charCodeAt(i) + ((hash << 5) - hash);
    hash = hash & hash; // Convert to 32-bit integer
  }

  // Use HSL for better color distribution
  // Hue: 0-360, Saturation: 40-60%, Lightness: 45-55%
  const hue = Math.abs(hash % 360);
  const saturation = 45 + (Math.abs(hash) % 15); // 45-60%
  const lightness = 45 + (Math.abs(hash >> 8) % 10); // 45-55%

  return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
}

/**
 * Alternative: Returns a Tailwind color class from the email
 * Useful for use with Tailwind background colors
 */
export function colorClassFromEmail(email: string): string {
  const colors = [
    "bg-red-500",
    "bg-orange-500",
    "bg-amber-500",
    "bg-yellow-500",
    "bg-lime-500",
    "bg-green-500",
    "bg-emerald-500",
    "bg-teal-500",
    "bg-cyan-500",
    "bg-sky-500",
    "bg-blue-500",
    "bg-indigo-500",
    "bg-violet-500",
    "bg-purple-500",
    "bg-fuchsia-500",
    "bg-pink-500",
    "bg-rose-500",
  ];

  if (!email) {
    return colors[0];
  }

  let hash = 0;
  for (let i = 0; i < email.length; i++) {
    hash = email.charCodeAt(i) + ((hash << 5) - hash);
    hash = hash & hash;
  }

  return colors[Math.abs(hash) % colors.length];
}
