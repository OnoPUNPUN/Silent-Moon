import 'package:flutter/widgets.dart';

/// Semantic color palette for the **Silent Moon** Personal Reading Library application.
///
/// Organized cleanly by domain:
/// - Core Brand & UI surfaces (Light, Dark / Silent Moon Night, Sepia parchment)
/// - Text Highlighting & Annotations (Yellow, Pink, Green, Blue, Purple, Orange)
/// - Reading Statuses (Want to Read, Currently Reading, Completed, Archived)
/// - Document Formats (PDF, EPUB)
/// - Sync & Connectivity Statuses (Synced, Syncing, Offline/Pending, Conflict/Error)
/// - Reading Progress & Reader Controls
abstract final class AppColors {
  // ---------------------------------------------------------------------------
  // 1. Brand Identity & Accents
  // ---------------------------------------------------------------------------
  /// Deep tranquil indigo representing the "Silent Moon" night sky.
  static const Color brandPrimary = Color(0xFF6366F1); // Indigo-500
  static const Color brandPrimaryLight = Color(0xFF4F46E5); // Indigo-600
  static const Color brandPrimaryDark = Color(0xFF818CF8); // Indigo-400
  static const Color brandMoonlight = Color(0xFFC7D2FE); // Indigo-200

  // ---------------------------------------------------------------------------
  // 2. Light Theme Colors (Clean Paper / Editorial)
  // ---------------------------------------------------------------------------
  static const Color lightBackground = Color(0xFFFFFFFF);
  static const Color lightForeground = Color(0xFF0F172A); // Slate-900
  static const Color lightCard = Color(0xFFFFFFFF);
  static const Color lightCardForeground = Color(0xFF0F172A);
  static const Color lightPopover = Color(0xFFFFFFFF);
  static const Color lightPopoverForeground = Color(0xFF0F172A);
  static const Color lightPrimary = Color(0xFF0F172A);
  static const Color lightPrimaryForeground = Color(0xFFF8FAFC);
  static const Color lightSecondary = Color(0xFFF1F5F9); // Slate-100
  static const Color lightSecondaryForeground = Color(0xFF0F172A);
  static const Color lightMuted = Color(0xFFF1F5F9);
  static const Color lightMutedForeground = Color(0xFF64748B); // Slate-500
  static const Color lightAccent = Color(0xFFF1F5F9);
  static const Color lightAccentForeground = Color(0xFF0F172A);
  static const Color lightDestructive = Color(0xFFEF4444); // Red-500
  static const Color lightDestructiveForeground = Color(0xFFF8FAFC);
  static const Color lightBorder = Color(0xFFE2E8F0); // Slate-200
  static const Color lightInput = Color(0xFFE2E8F0);
  static const Color lightRing = Color(0xFF6366F1); // Brand Indigo focus
  static const Color lightSelection = Color(0xFFBFDBFE); // Soft Blue-200

  // ---------------------------------------------------------------------------
  // 3. Dark Theme Colors ("Silent Moon" Night / Deep Slate)
  // ---------------------------------------------------------------------------
  static const Color darkBackground = Color(0xFF090D16); // Deep midnight slate
  static const Color darkForeground = Color(0xFFF8FAFC); // Slate-50
  static const Color darkCard = Color(0xFF0F172A); // Slate-900 elevated surface
  static const Color darkCardForeground = Color(0xFFF8FAFC);
  static const Color darkPopover = Color(0xFF0F172A);
  static const Color darkPopoverForeground = Color(0xFFF8FAFC);
  static const Color darkPrimary = Color(0xFFF8FAFC);
  static const Color darkPrimaryForeground = Color(0xFF090D16);
  static const Color darkSecondary = Color(0xFF1E293B); // Slate-800
  static const Color darkSecondaryForeground = Color(0xFFF8FAFC);
  static const Color darkMuted = Color(0xFF1E293B);
  static const Color darkMutedForeground = Color(0xFF94A3B8); // Slate-400
  static const Color darkAccent = Color(0xFF1E293B);
  static const Color darkAccentForeground = Color(0xFFF8FAFC);
  static const Color darkDestructive = Color(0xFFEF4444);
  static const Color darkDestructiveForeground = Color(0xFFF8FAFC);
  static const Color darkBorder = Color(0xFF1E293B);
  static const Color darkInput = Color(0xFF1E293B);
  static const Color darkRing = Color(0xFF818CF8); // Indigo-400
  static const Color darkSelection = Color(0xFF1E3A5F);

  // ---------------------------------------------------------------------------
  // 4. Sepia Theme Colors (Warm Vintage Parchment for Long Reading)
  // ---------------------------------------------------------------------------
  static const Color sepiaBackground = Color(0xFFFBF0D9); // Warm parchment
  static const Color sepiaForeground = Color(0xFF433422); // Deep espresso ink
  static const Color sepiaCard = Color(0xFFF3E7CA);
  static const Color sepiaCardForeground = Color(0xFF433422);
  static const Color sepiaPopover = Color(0xFFF3E7CA);
  static const Color sepiaPopoverForeground = Color(0xFF433422);
  static const Color sepiaPrimary = Color(0xFF5C4033); // Vintage brown
  static const Color sepiaPrimaryForeground = Color(0xFFFDFBF7);
  static const Color sepiaSecondary = Color(0xFFEADCC4);
  static const Color sepiaSecondaryForeground = Color(0xFF433422);
  static const Color sepiaMuted = Color(0xFFEADCC4);
  static const Color sepiaMutedForeground = Color(0xFF7D6551);
  static const Color sepiaAccent = Color(0xFFE2D1B3);
  static const Color sepiaAccentForeground = Color(0xFF433422);
  static const Color sepiaDestructive = Color(0xFFDC2626);
  static const Color sepiaDestructiveForeground = Color(0xFFFFFFFF);
  static const Color sepiaBorder = Color(0xFFDDD0B8);
  static const Color sepiaInput = Color(0xFFDDD0B8);
  static const Color sepiaRing = Color(0xFF9C7A54);
  static const Color sepiaSelection = Color(0xFFE8D5B5);

  // ---------------------------------------------------------------------------
  // 5. Annotation & Highlight Colors
  // Aligned with the API specification and REST endpoints:
  // Yellow (#FFF3A3), Pink (#FFC0CB), Green (#90EE90), Blue, Purple, Orange.
  // ---------------------------------------------------------------------------
  // Hex strings matching backend REST schema
  static const String highlightYellowHex = '#FFF3A3';
  static const String highlightPinkHex = '#FFC0CB';
  static const String highlightGreenHex = '#90EE90';
  static const String highlightBlueHex = '#BAE6FD';
  static const String highlightPurpleHex = '#E9D5FF';
  static const String highlightOrangeHex = '#FED7AA';

  // Light / Default Highlight Colors
  static const Color highlightYellow = Color(0xFFFFF3A3); // Warm yellow
  static const Color highlightPink = Color(0xFFFFC0CB); // Classic pink
  static const Color highlightGreen = Color(0xFF90EE90); // Soft mint green
  static const Color highlightBlue = Color(0xFFBAE6FD); // Sky blue
  static const Color highlightPurple = Color(0xFFE9D5FF); // Lavender
  static const Color highlightOrange = Color(0xFFFED7AA); // Peach orange

  // Dark Mode Semi-Transparent Highlight Overlays (keeps text readable at night)
  static const Color highlightYellowDark = Color(0x66FDE047);
  static const Color highlightPinkDark = Color(0x66F472B6);
  static const Color highlightGreenDark = Color(0x664ADE80);
  static const Color highlightBlueDark = Color(0x6638BDF8);
  static const Color highlightPurpleDark = Color(0x66C084FC);
  static const Color highlightOrangeDark = Color(0x66FB923C);

  // Annotation Types
  static const Color annotationUnderline = Color(0xFF3B82F6); // Blue underline
  static const Color annotationStrikethrough = Color(0xFFEF4444); // Red strike
  static const Color bookmarkColor = Color(0xFFE11D48); // Rose red ribbon
  static const Color noteStickyColor = Color(0xFFFEF08A); // Yellow note badge
  static const Color noteStickyDark = Color(0xFF713F12);

  // ---------------------------------------------------------------------------
  // 6. Reading Status Colors
  // - Want to Read
  // - Currently Reading
  // - Completed
  // - Archived
  // ---------------------------------------------------------------------------
  static const Color statusWantToRead = Color(0xFF3B82F6); // Blue-500
  static const Color statusWantToReadBg = Color(0xFFEFF6FF); // Blue-50
  static const Color statusWantToReadDarkBg = Color(0xFF1E3A5F);

  static const Color statusCurrentlyReading = Color(0xFFF59E0B); // Amber-500
  static const Color statusCurrentlyReadingBg = Color(0xFFFFFBEB); // Amber-50
  static const Color statusCurrentlyReadingDarkBg = Color(0xFF451A03);

  static const Color statusCompleted = Color(0xFF10B981); // Emerald-500
  static const Color statusCompletedBg = Color(0xFFECFDF5); // Emerald-50
  static const Color statusCompletedDarkBg = Color(0xFF064E3B);

  static const Color statusArchived = Color(0xFF64748B); // Slate-500
  static const Color statusArchivedBg = Color(0xFFF1F5F9); // Slate-100
  static const Color statusArchivedDarkBg = Color(0xFF1E293B);

  // ---------------------------------------------------------------------------
  // 7. Book Format Colors (PDF vs EPUB)
  // ---------------------------------------------------------------------------
  static const Color formatPdf = Color(0xFFEF4444); // PDF Red
  static const Color formatPdfBg = Color(0xFFFEE2E2);
  static const Color formatPdfDarkBg = Color(0xFF450A0A);

  static const Color formatEpub = Color(0xFF0EA5E9); // EPUB Cyan/Sky
  static const Color formatEpubBg = Color(0xFFE0F2FE);
  static const Color formatEpubDarkBg = Color(0xFF082F49);

  // ---------------------------------------------------------------------------
  // 8. Offline & Sync Status Colors
  // - Synced: Up to date with backend
  // - Syncing: Active upload/download in progress
  // - Pending / Offline: Changes queued locally
  // - Conflict / Error: Requires attention
  // ---------------------------------------------------------------------------
  static const Color syncSuccess = Color(0xFF10B981); // Emerald-500
  static const Color syncSuccessBg = Color(0xFFD1FAE5);

  static const Color syncActive = Color(0xFF3B82F6); // Blue-500
  static const Color syncActiveBg = Color(0xFFDBEAFE);

  static const Color syncPending = Color(0xFFF59E0B); // Amber-500
  static const Color syncPendingBg = Color(0xFFFEF3C7);

  static const Color syncError = Color(0xFFEF4444); // Red-500
  static const Color syncErrorBg = Color(0xFFFEE2E2);

  // ---------------------------------------------------------------------------
  // 9. Reading Progress Indicator Colors
  // ---------------------------------------------------------------------------
  static const Color progressIndicator = Color(0xFF6366F1); // Silent Moon Indigo
  static const Color progressTrackLight = Color(0xFFE2E8F0);
  static const Color progressTrackDark = Color(0xFF1E293B);
  static const Color progressTrackSepia = Color(0xFFE5D5BA);
}

/// Helper model representing a selectable annotation highlight color.
///
/// Convenient for UI pickers, toolbar badges, and serializing to/from hex
/// for backend REST sync.
@immutable
class AppHighlightColor {
  const AppHighlightColor({
    required this.name,
    required this.hex,
    required this.color,
    required this.darkColor,
    required this.textColor,
  });

  /// Name of the highlight color (e.g. 'Yellow', 'Pink').
  final String name;

  /// Hex representation stored in database and REST API (e.g. '#FFF3A3').
  final String hex;

  /// Light theme color representation.
  final Color color;

  /// Dark theme color representation.
  final Color darkColor;

  /// High-contrast text color when displayed inside a chip or badge.
  final Color textColor;

  /// Resolves the color for the given brightness.
  Color resolve(Brightness brightness) {
    return brightness == Brightness.dark ? darkColor : color;
  }

  /// Preset list of predefined highlight colors per the application specification.
  static const List<AppHighlightColor> presets = [
    AppHighlightColor(
      name: 'Yellow',
      hex: AppColors.highlightYellowHex,
      color: AppColors.highlightYellow,
      darkColor: AppColors.highlightYellowDark,
      textColor: Color(0xFF713F12),
    ),
    AppHighlightColor(
      name: 'Pink',
      hex: AppColors.highlightPinkHex,
      color: AppColors.highlightPink,
      darkColor: AppColors.highlightPinkDark,
      textColor: Color(0xFF831843),
    ),
    AppHighlightColor(
      name: 'Green',
      hex: AppColors.highlightGreenHex,
      color: AppColors.highlightGreen,
      darkColor: AppColors.highlightGreenDark,
      textColor: Color(0xFF14532D),
    ),
    AppHighlightColor(
      name: 'Blue',
      hex: AppColors.highlightBlueHex,
      color: AppColors.highlightBlue,
      darkColor: AppColors.highlightBlueDark,
      textColor: Color(0xFF0C4A6E),
    ),
    AppHighlightColor(
      name: 'Purple',
      hex: AppColors.highlightPurpleHex,
      color: AppColors.highlightPurple,
      darkColor: AppColors.highlightPurpleDark,
      textColor: Color(0xFF581C87),
    ),
    AppHighlightColor(
      name: 'Orange',
      hex: AppColors.highlightOrangeHex,
      color: AppColors.highlightOrange,
      darkColor: AppColors.highlightOrangeDark,
      textColor: Color(0xFF7C2D12),
    ),
  ];

  /// Look up an [AppHighlightColor] by hex string, falling back to Yellow.
  static AppHighlightColor fromHex(String? hex) {
    if (hex == null) return presets.first;
    final normalized = hex.toUpperCase();
    return presets.firstWhere(
      (h) => h.hex.toUpperCase() == normalized,
      orElse: () => presets.first,
    );
  }
}
