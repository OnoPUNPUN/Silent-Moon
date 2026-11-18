import 'package:flutter/material.dart';
import 'package:shadcn_ui/shadcn_ui.dart';

import 'app_colors.dart';

/// Central theme configuration for the **Silent Moon** application using `shadcn_ui`.
///
/// Provides:
/// - Light, Dark (Silent Moon Night), and Sepia reading themes
/// - Deep integration with Shadcn UI component themes (Buttons, Badges, Cards,
///   Dialogs, Sheets, Inputs, Progress, Tabs, Popovers, etc.)
/// - Custom domain colors mapped through [ShadColorScheme.custom]
/// - Material 3 interoperability for cross-platform hybrid layouts
abstract final class AppTheme {
  AppTheme._();

  // ---------------------------------------------------------------------------
  // Global Metrics
  // ---------------------------------------------------------------------------
  static const BorderRadius globalRadius = BorderRadius.all(Radius.circular(8));
  static const BorderRadius cardRadius = BorderRadius.all(Radius.circular(10));
  static const BorderRadius modalRadius = BorderRadius.all(Radius.circular(12));
  static const BorderRadius badgeRadius = BorderRadius.all(Radius.circular(6));

  // ---------------------------------------------------------------------------
  // 1. Color Schemes
  // ---------------------------------------------------------------------------

  /// Light color scheme with Silent Moon domain-specific custom colors.
  static final ShadColorScheme lightColorScheme = ShadSlateColorScheme.light(
    background: AppColors.lightBackground,
    foreground: AppColors.lightForeground,
    card: AppColors.lightCard,
    cardForeground: AppColors.lightCardForeground,
    popover: AppColors.lightPopover,
    popoverForeground: AppColors.lightPopoverForeground,
    primary: AppColors.lightPrimary,
    primaryForeground: AppColors.lightPrimaryForeground,
    secondary: AppColors.lightSecondary,
    secondaryForeground: AppColors.lightSecondaryForeground,
    muted: AppColors.lightMuted,
    mutedForeground: AppColors.lightMutedForeground,
    accent: AppColors.lightAccent,
    accentForeground: AppColors.lightAccentForeground,
    destructive: AppColors.lightDestructive,
    destructiveForeground: AppColors.lightDestructiveForeground,
    border: AppColors.lightBorder,
    input: AppColors.lightInput,
    ring: AppColors.lightRing,
    selection: AppColors.lightSelection,
    custom: _buildCustomColorsMap(brightness: Brightness.light),
  );

  /// Dark color scheme tuned for midnight reading ("Silent Moon" aesthetic).
  static final ShadColorScheme darkColorScheme = ShadSlateColorScheme.dark(
    background: AppColors.darkBackground,
    foreground: AppColors.darkForeground,
    card: AppColors.darkCard,
    cardForeground: AppColors.darkCardForeground,
    popover: AppColors.darkPopover,
    popoverForeground: AppColors.darkPopoverForeground,
    primary: AppColors.darkPrimary,
    primaryForeground: AppColors.darkPrimaryForeground,
    secondary: AppColors.darkSecondary,
    secondaryForeground: AppColors.darkSecondaryForeground,
    muted: AppColors.darkMuted,
    mutedForeground: AppColors.darkMutedForeground,
    accent: AppColors.darkAccent,
    accentForeground: AppColors.darkAccentForeground,
    destructive: AppColors.darkDestructive,
    destructiveForeground: AppColors.darkDestructiveForeground,
    border: AppColors.darkBorder,
    input: AppColors.darkInput,
    ring: AppColors.darkRing,
    selection: AppColors.darkSelection,
    custom: _buildCustomColorsMap(brightness: Brightness.dark),
  );

  /// Sepia color scheme for warm, eye-friendly daytime reading.
  static final ShadColorScheme sepiaColorScheme = ShadColorScheme(
    background: AppColors.sepiaBackground,
    foreground: AppColors.sepiaForeground,
    card: AppColors.sepiaCard,
    cardForeground: AppColors.sepiaCardForeground,
    popover: AppColors.sepiaPopover,
    popoverForeground: AppColors.sepiaPopoverForeground,
    primary: AppColors.sepiaPrimary,
    primaryForeground: AppColors.sepiaPrimaryForeground,
    secondary: AppColors.sepiaSecondary,
    secondaryForeground: AppColors.sepiaSecondaryForeground,
    muted: AppColors.sepiaMuted,
    mutedForeground: AppColors.sepiaMutedForeground,
    accent: AppColors.sepiaAccent,
    accentForeground: AppColors.sepiaAccentForeground,
    destructive: AppColors.sepiaDestructive,
    destructiveForeground: AppColors.sepiaDestructiveForeground,
    border: AppColors.sepiaBorder,
    input: AppColors.sepiaInput,
    ring: AppColors.sepiaRing,
    selection: AppColors.sepiaSelection,
    custom: _buildCustomColorsMap(
      brightness: Brightness.light,
      isSepia: true,
    ),
  );

  // ---------------------------------------------------------------------------
  // 2. Shadcn Theme Definitions
  // ---------------------------------------------------------------------------

  /// Standard Light Theme for the application.
  static ShadThemeData get light => _buildTheme(
        brightness: Brightness.light,
        scheme: lightColorScheme,
      );

  /// Silent Moon Dark Theme for night-time reading.
  static ShadThemeData get dark => _buildTheme(
        brightness: Brightness.dark,
        scheme: darkColorScheme,
      );

  /// Dedicated Sepia Reading Theme for comfortable long-form reading.
  static ShadThemeData get sepia => _buildTheme(
        brightness: Brightness.light,
        scheme: sepiaColorScheme,
      );

  // ---------------------------------------------------------------------------
  // 3. Theme Builder
  // ---------------------------------------------------------------------------
  static ShadThemeData _buildTheme({
    required Brightness brightness,
    required ShadColorScheme scheme,
  }) {
    final isDark = brightness == Brightness.dark;

    return ShadThemeData(
      brightness: brightness,
      colorScheme: scheme,
      radius: globalRadius,

      // --- Buttons ---
      primaryButtonTheme: const ShadButtonTheme(
        size: ShadButtonSize.regular,
      ),
      secondaryButtonTheme: const ShadButtonTheme(
        size: ShadButtonSize.regular,
      ),
      outlineButtonTheme: const ShadButtonTheme(
        size: ShadButtonSize.regular,
      ),
      ghostButtonTheme: const ShadButtonTheme(
        size: ShadButtonSize.regular,
      ),
      destructiveButtonTheme: const ShadButtonTheme(
        size: ShadButtonSize.regular,
      ),

      // --- Badges (for book formats, reading status, and tag pills) ---
      primaryBadgeTheme: const ShadBadgeTheme(
        shape: RoundedRectangleBorder(borderRadius: badgeRadius),
        padding: EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      ),
      secondaryBadgeTheme: const ShadBadgeTheme(
        shape: RoundedRectangleBorder(borderRadius: badgeRadius),
        padding: EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      ),
      outlineBadgeTheme: const ShadBadgeTheme(
        shape: RoundedRectangleBorder(borderRadius: badgeRadius),
        padding: EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      ),
      destructiveBadgeTheme: const ShadBadgeTheme(
        shape: RoundedRectangleBorder(borderRadius: badgeRadius),
        padding: EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      ),

      // --- Cards (Books in library, collections, annotations) ---
      cardTheme: ShadCardTheme(
        radius: cardRadius,
        padding: const EdgeInsets.all(16),
        backgroundColor: scheme.card,
        border: ShadBorder.all(
          color: scheme.border,
          width: 1,
          radius: cardRadius,
        ),
        shadows: isDark
            ? const [
                BoxShadow(
                  color: Color(0x33000000),
                  blurRadius: 8,
                  offset: Offset(0, 2),
                ),
              ]
            : const [
                BoxShadow(
                  color: Color(0x0A000000),
                  blurRadius: 6,
                  offset: Offset(0, 2),
                ),
              ],
      ),

      // --- Dialogs (delete confirmation, sync conflict, edit tag) ---
      primaryDialogTheme: ShadDialogTheme(
        radius: modalRadius,
        backgroundColor: scheme.card,
        border: Border.all(color: scheme.border, width: 1),
        padding: const EdgeInsets.all(24),
      ),
      alertDialogTheme: ShadDialogTheme(
        radius: modalRadius,
        backgroundColor: scheme.card,
        border: Border.all(color: scheme.border, width: 1),
        padding: const EdgeInsets.all(24),
      ),

      // --- Sheets (table of contents, annotations panel, book details drawer) ---
      sheetTheme: ShadSheetTheme(
        radius: modalRadius,
        backgroundColor: scheme.card,
        border: Border.all(color: scheme.border, width: 1),
        padding: const EdgeInsets.all(20),
      ),

      // --- Popovers (highlight color picker, reading font settings) ---
      popoverTheme: ShadPopoverTheme(
        padding: const EdgeInsets.all(12),
        decoration: ShadDecoration(
          color: scheme.popover,
          border: ShadBorder.all(color: scheme.border, width: 1, radius: globalRadius),
          shadows: const [
            BoxShadow(
              color: Color(0x1F000000),
              blurRadius: 10,
              offset: Offset(0, 4),
            ),
          ],
        ),
      ),

      // --- Tooltips ---
      tooltipTheme: ShadTooltipTheme(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
        decoration: ShadDecoration(
          color: isDark ? const Color(0xFF1E293B) : const Color(0xFF0F172A),
          border: ShadBorder.all(
            color: isDark ? const Color(0xFF334155) : const Color(0xFF1E293B),
            width: 1,
            radius: const BorderRadius.all(Radius.circular(6)),
          ),
        ),
      ),

      // --- Inputs & Search Fields ---
      inputTheme: ShadInputTheme(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
        decoration: ShadDecoration(
          color: scheme.card,
          border: ShadBorder.all(color: scheme.border, width: 1, radius: globalRadius),
          focusedBorder: ShadBorder.all(color: scheme.ring, width: 2, radius: globalRadius),
        ),
      ),

      // --- Selects (collection selector, sort dropdown) ---
      selectTheme: ShadSelectTheme(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
        decoration: ShadDecoration(
          color: scheme.card,
          border: ShadBorder.all(color: scheme.border, width: 1, radius: globalRadius),
        ),
      ),

      // --- Reading Progress Indicator ---
      progressTheme: ShadProgressTheme(
        minHeight: 6,
        borderRadius: const BorderRadius.all(Radius.circular(999)),
        color: AppColors.progressIndicator,
        backgroundColor: isDark
            ? AppColors.progressTrackDark
            : (scheme.background == AppColors.sepiaBackground
                ? AppColors.progressTrackSepia
                : AppColors.progressTrackLight),
      ),

      // --- Navigation Tabs ---
      tabsTheme: const ShadTabsTheme(),

      // --- Toast Notifications (sync alerts, download complete) ---
      primaryToastTheme: ShadToastTheme(
        backgroundColor: scheme.card,
        border: ShadBorder.all(color: scheme.border, width: 1, radius: globalRadius),
      ),
      destructiveToastTheme: ShadToastTheme(
        backgroundColor: scheme.destructive,
        border: ShadBorder.all(color: scheme.destructive, width: 1, radius: globalRadius),
      ),

      // --- Separators & Dividers ---
      separatorTheme: ShadSeparatorTheme(
        color: scheme.border,
      ),

      // --- Context Menus (desktop right click, reader highlight actions) ---
      contextMenuTheme: const ShadContextMenuTheme(),
    );
  }

  // ---------------------------------------------------------------------------
  // 4. Material Theme Interoperability
  // ---------------------------------------------------------------------------

  /// Creates a Material [ThemeData] corresponding to [light].
  static ThemeData get materialLight => _buildMaterialTheme(
        brightness: Brightness.light,
        scheme: lightColorScheme,
      );

  /// Creates a Material [ThemeData] corresponding to [dark].
  static ThemeData get materialDark => _buildMaterialTheme(
        brightness: Brightness.dark,
        scheme: darkColorScheme,
      );

  /// Creates a Material [ThemeData] corresponding to [sepia].
  static ThemeData get materialSepia => _buildMaterialTheme(
        brightness: Brightness.light,
        scheme: sepiaColorScheme,
      );

  static ThemeData _buildMaterialTheme({
    required Brightness brightness,
    required ShadColorScheme scheme,
  }) {
    return ThemeData(
      brightness: brightness,
      scaffoldBackgroundColor: scheme.background,
      colorScheme: ColorScheme(
        brightness: brightness,
        primary: scheme.primary,
        onPrimary: scheme.primaryForeground,
        secondary: scheme.secondary,
        onSecondary: scheme.secondaryForeground,
        error: scheme.destructive,
        onError: scheme.destructiveForeground,
        surface: scheme.background,
        onSurface: scheme.foreground,
      ),
      dividerTheme: DividerThemeData(
        color: scheme.border,
        thickness: 1,
      ),
      textSelectionTheme: TextSelectionThemeData(
        cursorColor: scheme.primary,
        selectionColor: scheme.selection,
        selectionHandleColor: scheme.primary,
      ),
      iconTheme: IconThemeData(
        size: 18,
        color: scheme.foreground,
      ),
      cardTheme: CardThemeData(
        color: scheme.card,
        elevation: 0,
        shape: RoundedRectangleBorder(
          side: BorderSide(color: scheme.border),
          borderRadius: cardRadius,
        ),
      ),
    );
  }

  // ---------------------------------------------------------------------------
  // 5. Custom Colors Registry
  // ---------------------------------------------------------------------------
  static Map<String, Color> _buildCustomColorsMap({
    required Brightness brightness,
    bool isSepia = false,
  }) {
    final isDark = brightness == Brightness.dark;

    return {
      // Reading Statuses
      'statusWantToRead': AppColors.statusWantToRead,
      'statusWantToReadBg':
          isDark ? AppColors.statusWantToReadDarkBg : AppColors.statusWantToReadBg,
      'statusCurrentlyReading': AppColors.statusCurrentlyReading,
      'statusCurrentlyReadingBg':
          isDark ? AppColors.statusCurrentlyReadingDarkBg : AppColors.statusCurrentlyReadingBg,
      'statusCompleted': AppColors.statusCompleted,
      'statusCompletedBg':
          isDark ? AppColors.statusCompletedDarkBg : AppColors.statusCompletedBg,
      'statusArchived': AppColors.statusArchived,
      'statusArchivedBg':
          isDark ? AppColors.statusArchivedDarkBg : AppColors.statusArchivedBg,

      // Book Formats
      'formatPdf': AppColors.formatPdf,
      'formatPdfBg': isDark ? AppColors.formatPdfDarkBg : AppColors.formatPdfBg,
      'formatEpub': AppColors.formatEpub,
      'formatEpubBg': isDark ? AppColors.formatEpubDarkBg : AppColors.formatEpubBg,

      // Sync & Network
      'syncSuccess': AppColors.syncSuccess,
      'syncSuccessBg': AppColors.syncSuccessBg,
      'syncActive': AppColors.syncActive,
      'syncActiveBg': AppColors.syncActiveBg,
      'syncPending': AppColors.syncPending,
      'syncPendingBg': AppColors.syncPendingBg,
      'syncError': AppColors.syncError,
      'syncErrorBg': AppColors.syncErrorBg,

      // Highlights (per REST API schema)
      'highlightYellow': isDark ? AppColors.highlightYellowDark : AppColors.highlightYellow,
      'highlightPink': isDark ? AppColors.highlightPinkDark : AppColors.highlightPink,
      'highlightGreen': isDark ? AppColors.highlightGreenDark : AppColors.highlightGreen,
      'highlightBlue': isDark ? AppColors.highlightBlueDark : AppColors.highlightBlue,
      'highlightPurple': isDark ? AppColors.highlightPurpleDark : AppColors.highlightPurple,
      'highlightOrange': isDark ? AppColors.highlightOrangeDark : AppColors.highlightOrange,

      // Annotations & Reader Markers
      'bookmark': AppColors.bookmarkColor,
      'annotationUnderline': AppColors.annotationUnderline,
      'annotationStrikethrough': AppColors.annotationStrikethrough,
      'noteSticky': isDark ? AppColors.noteStickyDark : AppColors.noteStickyColor,

      // Progress Track
      'progressTrack': isDark
          ? AppColors.progressTrackDark
          : (isSepia ? AppColors.progressTrackSepia : AppColors.progressTrackLight),
      'progressIndicator': AppColors.progressIndicator,
    };
  }
}

// -----------------------------------------------------------------------------
// ShadColorScheme Custom Domain Accessors Extension
// -----------------------------------------------------------------------------

/// Extension providing type-safe access to domain-specific colors from [ShadColorScheme].
extension AppColorSchemeExtension on ShadColorScheme {
  // Reading Status Colors
  Color get statusWantToRead => custom['statusWantToRead'] ?? AppColors.statusWantToRead;
  Color get statusWantToReadBg => custom['statusWantToReadBg'] ?? AppColors.statusWantToReadBg;
  Color get statusCurrentlyReading =>
      custom['statusCurrentlyReading'] ?? AppColors.statusCurrentlyReading;
  Color get statusCurrentlyReadingBg =>
      custom['statusCurrentlyReadingBg'] ?? AppColors.statusCurrentlyReadingBg;
  Color get statusCompleted => custom['statusCompleted'] ?? AppColors.statusCompleted;
  Color get statusCompletedBg => custom['statusCompletedBg'] ?? AppColors.statusCompletedBg;
  Color get statusArchived => custom['statusArchived'] ?? AppColors.statusArchived;
  Color get statusArchivedBg => custom['statusArchivedBg'] ?? AppColors.statusArchivedBg;

  // Book Format Colors
  Color get formatPdf => custom['formatPdf'] ?? AppColors.formatPdf;
  Color get formatPdfBg => custom['formatPdfBg'] ?? AppColors.formatPdfBg;
  Color get formatEpub => custom['formatEpub'] ?? AppColors.formatEpub;
  Color get formatEpubBg => custom['formatEpubBg'] ?? AppColors.formatEpubBg;

  // Sync & Offline Status Colors
  Color get syncSuccess => custom['syncSuccess'] ?? AppColors.syncSuccess;
  Color get syncSuccessBg => custom['syncSuccessBg'] ?? AppColors.syncSuccessBg;
  Color get syncActive => custom['syncActive'] ?? AppColors.syncActive;
  Color get syncActiveBg => custom['syncActiveBg'] ?? AppColors.syncActiveBg;
  Color get syncPending => custom['syncPending'] ?? AppColors.syncPending;
  Color get syncPendingBg => custom['syncPendingBg'] ?? AppColors.syncPendingBg;
  Color get syncError => custom['syncError'] ?? AppColors.syncError;
  Color get syncErrorBg => custom['syncErrorBg'] ?? AppColors.syncErrorBg;

  // Highlights
  Color get highlightYellow => custom['highlightYellow'] ?? AppColors.highlightYellow;
  Color get highlightPink => custom['highlightPink'] ?? AppColors.highlightPink;
  Color get highlightGreen => custom['highlightGreen'] ?? AppColors.highlightGreen;
  Color get highlightBlue => custom['highlightBlue'] ?? AppColors.highlightBlue;
  Color get highlightPurple => custom['highlightPurple'] ?? AppColors.highlightPurple;
  Color get highlightOrange => custom['highlightOrange'] ?? AppColors.highlightOrange;

  // Annotations & Markers
  Color get bookmark => custom['bookmark'] ?? AppColors.bookmarkColor;
  Color get annotationUnderline =>
      custom['annotationUnderline'] ?? AppColors.annotationUnderline;
  Color get annotationStrikethrough =>
      custom['annotationStrikethrough'] ?? AppColors.annotationStrikethrough;
  Color get noteSticky => custom['noteSticky'] ?? AppColors.noteStickyColor;

  // Progress
  Color get progressTrack => custom['progressTrack'] ?? AppColors.progressTrackLight;
  Color get progressIndicator => custom['progressIndicator'] ?? AppColors.progressIndicator;
}

// -----------------------------------------------------------------------------
// BuildContext Convenience Extensions
// -----------------------------------------------------------------------------

/// Context extension providing convenient shortcuts to theme properties.
extension AppThemeContextExtension on BuildContext {
  /// The active [ShadThemeData].
  ShadThemeData get shadTheme => ShadTheme.of(this);

  /// The active [ShadColorScheme].
  ShadColorScheme get colorScheme => ShadTheme.of(this).colorScheme;

  /// The active [ShadTextTheme].
  ShadTextTheme get textTheme => ShadTheme.of(this).textTheme;

  /// Whether the active theme brightness is dark.
  bool get isDarkMode => ShadTheme.of(this).brightness == Brightness.dark;
}
