import 'package:flutter/material.dart';
import 'package:frontend/models/create_alert.model.dart';
import 'package:frontend/models/response_model.dart';
import 'package:frontend/services/alert/alert_api_impl.dart';

class CreateAlertController extends ChangeNotifier {
  final AlertApiImpl alertApiService;

  CreateAlertController({required this.alertApiService});
  final formKey = GlobalKey<FormState>();
  final ScrollController scrollController = ScrollController();
  final TextEditingController alertNameController = TextEditingController();
  final TextEditingController ttmController = TextEditingController();
  final TextEditingController currentPriceController = TextEditingController();
  final TextEditingController alertPriceController = TextEditingController();

  Map<String, String?> errorText = {
    "alertName": null,
    "tickerToMonitor": null,
    "currentPrice": null,
    "alertPrice": null,
  };
  Map<String, String?> fieldValue = {
    "alertName": null,
    "tickerToMonitor": null,
    "condition": null,
    "alertPrice": null,
    "currentFetchedPrice": null,
    "currentFetchedTime": null
  };
  Map<String, String> condiitionTextMap = {
    ">": "goes above",
    ">=": "is equal or goes above",
    "==": "is equal",
    "<": "goes under",
    "<=": "is equal or goes under",
  };
  bool isPasswordVisible = false;
  bool isLoading = false;
  bool isGetCurrentPriceLoading = false;
  String? errorMessage;
  String? successMessage;

  void validateField(String field, String? errorMsg) {
    errorText[field] = errorMsg;
    notifyListeners();
  }

  void setFieldValue(String field, String? value) {
    fieldValue[field] = value;
    notifyListeners();
  }

  Future<void> createAlert() async {
    if (!formKey.currentState!.validate()) return;

    isLoading = true;
    errorMessage = null;
    successMessage = null;
    notifyListeners();

    try {
      ResponseModel<CreateAlertModel> result =
          await alertApiService.createAlert(
        CreateAlertModel(
          alertName: alertNameController.text,
          tickerToMonitor: ttmController.text,
          currentFetchedPrice: double.parse(currentPriceController.text),
          condition: fieldValue["condition"]!,
          alertPrice: double.parse(alertPriceController.text),
          currentFetchedTime: fieldValue["currentFetchedTime"]!,
        ),
      );
      if (result.success) {
        successMessage = result.message;
      } else {
        errorMessage = result.message;
      }
    } catch (e) {
      errorMessage = "Create alert failed: $e";
    } finally {
      isLoading = false;
      notifyListeners();
    }
    @override
    void dispose() {
      // Clean up the controller when the widget is disposed.
      alertNameController.dispose();
      ttmController.dispose();
      currentPriceController.dispose();
      alertPriceController.dispose();
      super.dispose();
    }
  }

  Future<void> getcurrentPriceAndTime() async {
    isGetCurrentPriceLoading = true;
    errorMessage = null;
    successMessage = null;
    notifyListeners();

    try {
      ResponseModel<GetCurrentPriceAndTimeModel> result =
          await alertApiService.getcurrentPriceAndTime(
        GetCurrentPriceAndTimeModel(
            tickerToMonitor: "${ttmController.text}.NS"),
      );

      if (result.success && result.data != null) {
        debugPrint("API Response successful");
        debugPrint("Parsed Data: ${result.data.toString()}");

        // Update UI fields
        if (result.data!.currentFetchedPrice != null) {
          currentPriceController.text =
              result.data!.currentFetchedPrice.toString();
          fieldValue["currentFetchedTime"] = result.data!.currentFetchedTime;
          successMessage =
              "Price fetched successfully: ₹${result.data!.currentFetchedPrice}";
        } else {
          errorMessage = "Failed to parse price data";
        }

        notifyListeners();
      } else {
        errorMessage = result.message;
        currentPriceController.text = "0";
        debugPrint("API Error: $errorMessage");
      }
    } catch (e) {
      errorMessage = "Failed to fetch price: $e";
      debugPrint("Exception: $e");
    } finally {
      isGetCurrentPriceLoading = false;
      notifyListeners();
    }
  }
}
