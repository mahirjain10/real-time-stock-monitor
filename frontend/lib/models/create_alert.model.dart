import 'package:intl/intl.dart';
import 'dart:convert';

class CreateAlertModel {
  // final String userId;
  final String alertName;
  final String tickerToMonitor;
  final double currentFetchedPrice;
  final String condition;
  final double alertPrice;
  final String? currentFetchedTime; // Assuming it represents a timestamp

  CreateAlertModel({
    // required this.userId,
    required this.alertName,
    required this.tickerToMonitor,
    required this.currentFetchedPrice,
    required this.condition,
    required this.alertPrice,
    required this.currentFetchedTime,
  });

  /// Convert object to JSON
  Map<String, dynamic> toJson() {
    return {
      // "user_id": userId,
      "alert_name": alertName,
      "ticker_to_monitor": tickerToMonitor,
      "current_fetched_price": currentFetchedPrice,
      "alert_condition": condition,
      "alert_price": alertPrice,
      "current_fetched_time": currentFetchedTime, // Convert DateTime to String
    };
  }

  /// Convert object to JSON string
  String toJsonString() {
    return jsonEncode(toJson());
  }

  /// Create object from JSON map
  factory CreateAlertModel.fromJson(Map<String, dynamic> json) {
    return CreateAlertModel(
        // userId: json["user_id"] as String,
        alertName: json["alert_name"] as String,
        tickerToMonitor: json["ticker_to_monitor"] as String,
        currentFetchedPrice: (json["current_fetched_price"] as num).toDouble(),
        // targetPrice: (json["target_price"] as num).toDouble(),
        condition: json["alert_condition"] as String,
        alertPrice: (json["alert_price"] as num).toDouble(),
        currentFetchedTime: json["current_fetched_time"]
        // ? DateTime.tryParse(json["current_fetched_time"])
        // : null,
        );
  }

  /// Create object from JSON string
  factory CreateAlertModel.fromJsonString(String jsonString) {
    return CreateAlertModel.fromJson(jsonDecode(jsonString));
  }
}

// import 'package:intl/intl.dart';

class GetCurrentPriceAndTimeModel {
  final String? tickerToMonitor;
  final double? currentFetchedPrice;
  final String? currentFetchedTime;

  GetCurrentPriceAndTimeModel({
    this.tickerToMonitor,
    this.currentFetchedPrice,
    this.currentFetchedTime,
  });

  /// Convert object to JSON for sending request
  Map<String, dynamic> toJson() {
    return {
      "ticker_to_monitor": tickerToMonitor,
    };
  }

  /// Convert object to JSON string
  String toJsonString() {
    return jsonEncode(toJson());
  }

  factory GetCurrentPriceAndTimeModel.fromJson(Map<String, dynamic> json) {
    // Handle the nested data structure
    final nestedData = json['data'] as Map<String, dynamic>;
    return GetCurrentPriceAndTimeModel(
      currentFetchedPrice:
          (nestedData['current_fetched_price'] as num?)?.toDouble(),
      currentFetchedTime: nestedData['current_fetched_time'] as String?,
    );
  }

  @override
  String toString() {
    return 'GetCurrentPriceAndTimeModel(tickerToMonitor: $tickerToMonitor, '
        'currentFetchedPrice: $currentFetchedPrice, currentFetchedTime: $currentFetchedTime)';
  }
}
