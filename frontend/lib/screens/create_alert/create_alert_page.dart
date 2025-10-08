import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:frontend/main.dart';
import 'package:frontend/screens/create_alert/alert_label_inputbar.dart';
import 'package:frontend/screens/create_alert/alert_preview.dart';
import 'package:frontend/screens/create_alert/controller/create_alert_controller.dart';
import 'package:frontend/screens/create_alert/create_alert_app_bar.dart';
import 'package:frontend/services/alert/alert_api_impl.dart';
import 'package:frontend/widgets/segment_slider.dart';
import 'package:provider/provider.dart';
import 'package:toastification/toastification.dart';
// import 'package:segmented_button_slide/segmented_button_slide.dart';

class CreateAlertPage extends StatefulWidget {
  const CreateAlertPage({super.key});

  @override
  State<CreateAlertPage> createState() => _CreateAlertPageState();
}

class _CreateAlertPageState extends State<CreateAlertPage> {
  @override
  Widget build(BuildContext context) {
    final screenWidth = MediaQuery.of(context).size.width;
    final screenHeight = MediaQuery.of(context).size.height;

    void _showToast(BuildContext context, String message, bool isError) {
      debugPrint(message);
      toastification.show(
        context: context,
        type: isError ? ToastificationType.error : ToastificationType.success,
        style: ToastificationStyle.flat,
        title: Text(isError ? "Error" : "Success",
            style: const TextStyle(fontWeight: FontWeight.bold)),
        description: Text(message),
        alignment: Alignment.topCenter,
        autoCloseDuration: const Duration(seconds: 4),
      );
    }

    return SafeArea(
      child: Scaffold(
        backgroundColor: Colors.grey.shade100,
        appBar: CreateAlertAppBar(),
        body: ChangeNotifierProvider(
          create: (context) =>
              CreateAlertController(alertApiService: AlertApiImpl(dio)),
          child: Consumer<CreateAlertController>(
            builder: (context, createAlertController, child) {
              // Show toasts when success or error messages change
              WidgetsBinding.instance.addPostFrameCallback((_) {
                if (createAlertController.successMessage != null) {
                  _showToast(
                      context, createAlertController.successMessage!, false);
                  createAlertController.successMessage =
                      null; // Reset after showing
                }
                if (createAlertController.errorMessage != null) {
                  _showToast(
                      context, createAlertController.errorMessage!, true);
                  createAlertController.errorMessage =
                      null; // Reset after showing
                }
              });

              return Stack(
                children: [
                  Column(
                    children: [
                      const SizedBox(height: 10),
                      SizedBox(
                        // height: screenHeight * 0.18,
                        width: screenWidth * 0.95,
                        child: AlertPreview(
                          createAlertController,
                        ),
                        // child: AlertPreview(controller: createAlertController),
                      ),
                      Expanded(
                        child: SingleChildScrollView(
                          controller: createAlertController.scrollController,
                          padding: const EdgeInsets.symmetric(horizontal: 16.0),
                          child: Form(
                            key: createAlertController.formKey,
                            child: Column(
                              children: [
                                const SizedBox(height: 20),
                                AlertLabelInputbar(
                                  controller:
                                      createAlertController.alertNameController,
                                  field: "alertName",
                                  validateForm:
                                      createAlertController.validateField,
                                  label: "Alert Name",
                                  hintText: "Enter alert name",
                                  errorText: createAlertController
                                      .errorText["alertName"],
                                  setFieldValue:
                                      createAlertController.setFieldValue,
                                ),
                                AlertLabelInputbar(
                                  controller:
                                      createAlertController.ttmController,
                                  field: "tickerToMonitor",
                                  validateForm:
                                      createAlertController.validateField,
                                  label: "Ticker to monitor",
                                  hintText: "Enter stock ticker (e.g., AAPL)",
                                  errorText: createAlertController
                                      .errorText["tickerToMonitor"],
                                  setFieldValue:
                                      createAlertController.setFieldValue,
                                  getCurrentPrice: createAlertController
                                      .getcurrentPriceAndTime,
                                ),
                                AlertLabelInputbar(
                                  controller: createAlertController
                                      .currentPriceController,
                                  field: "currentPrice",
                                  validateForm:
                                      createAlertController.validateField,
                                  label: "Current Price",
                                  hintText: "\u{20B9} Enter current price",
                                  errorText: createAlertController
                                      .errorText["currentPrice"],
                                  setFieldValue:
                                      createAlertController.setFieldValue,
                                ),
                                SegmentedButtonSlideWidget(
                                  setFieldValue:
                                      createAlertController.setFieldValue,
                                ),
                                AlertLabelInputbar(
                                  controller: createAlertController
                                      .alertPriceController,
                                  field: "alertPrice",
                                  validateForm:
                                      createAlertController.validateField,
                                  label: "Enter target price",
                                  hintText: "\u{20B9} Enter target price",
                                  errorText: createAlertController
                                      .errorText["alertPrice"],
                                  setFieldValue:
                                      createAlertController.setFieldValue,
                                ),
                                const SizedBox(height: 20),
                                SizedBox(
                                  height: screenHeight * 0.06,
                                  width: screenWidth * 0.9,
                                  child: ElevatedButton(
                                    style: ButtonStyle(
                                      backgroundColor: WidgetStateProperty.all(
                                          const Color(0xFF1D4ED8)),
                                      shape: WidgetStateProperty.all<
                                          RoundedRectangleBorder>(
                                        RoundedRectangleBorder(
                                          borderRadius:
                                              BorderRadius.circular(5.0),
                                        ),
                                      ),
                                    ),
                                    onPressed: createAlertController.isLoading
                                        ? null
                                        : () async {
                                            await createAlertController
                                                .createAlert();
                                            Future.microtask(() {
                                              if (createAlertController
                                                      .errorMessage !=
                                                  null) {
                                                _showToast(
                                                    context,
                                                    createAlertController
                                                        .errorMessage!,
                                                    true);
                                              } else if (createAlertController
                                                      .successMessage !=
                                                  null) {
                                                _showToast(
                                                    context,
                                                    createAlertController
                                                        .successMessage!,
                                                    false);
                                              }
                                            });
                                          },
                                    child: createAlertController
                                            .isGetCurrentPriceLoading
                                        ? const CircularProgressIndicator(
                                            color: Colors.white)
                                        : const Text(
                                            "Create Alert",
                                            style: TextStyle(
                                                color: Colors.white,
                                                fontSize: 20,
                                                fontWeight: FontWeight.w200),
                                          ),
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                  if (createAlertController.isGetCurrentPriceLoading)
                    Container(
                      color: Colors.black.withOpacity(
                          0.3), // Optional: Adds a semi-transparent overlay
                      child: const Center(
                        child: CircularProgressIndicator(),
                      ),
                    ),
                ],
              );
            },
          ),
        ),
      ),
    );
  }
}
