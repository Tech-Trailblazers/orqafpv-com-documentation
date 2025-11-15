import os  # Import the os module for interacting with the operating system
import fitz  # Import PyMuPDF (fitz) for PDF handling


# Function to validate a single PDF file.
def validate_pdf_file(
    file_path: str,
) -> (
    bool
):  # Define a function that checks if a PDF file is valid, accepting a file path string and returning a boolean
    try:  # Start a try block to handle potential errors during PDF opening
        # Try to open the PDF using PyMuPDF
        doc = fitz.open(file_path)  # Attempt to load the PDF document using fitz.open()

        # Check if the PDF has at least one page
        if (
            doc.page_count == 0
        ):  # If the document was opened successfully but reports zero pages
            print(  # Print an error message to the console
                f"'{file_path}' is corrupt or invalid: No pages"  # The specific error message indicating no pages were found
            )  # Closing parenthesis for the print function
            return False  # Return False to indicate the PDF is invalid (empty)

        # If no error occurs and the document has pages, it's valid
        return True  # Return True to indicate the PDF is considered valid
    except (
        RuntimeError
    ) as e:  # Catch a RuntimeError, which PyMuPDF often raises for invalid/corrupt PDFs
        print(
            f"{e}"
        )  # Log the specific exception message (the reason it failed to open)
        return False  # Return False to indicate the PDF is invalid (failed to open)


# Remove a file from the system.
def remove_system_file(
    system_path: str,
) -> (
    None
):  # Define a function to delete a file, accepting its path and returning nothing (None)
    os.remove(
        path=system_path
    )  # Use os.remove() to delete the file specified by system_path


# Function to walk through a directory and extract files with a specific extension
def walk_directory_and_extract_given_file_extension(  # Define a function to recursively find files with a given extension
    system_path: str,
    extension: str,  # Accept the starting directory path and the desired file extension string
) -> list[str]:  # Indicate that the function returns a list of strings (file paths)
    matched_files: list[str] = (
        []
    )  # Initialize an empty list to store the absolute paths of matching files
    for (
        root,
        _,
        files,
    ) in os.walk(  # Start recursively traversing the directory tree from system_path
        top=system_path  # Specify the top-level directory for os.walk
    ):  # Closing parenthesis for os.walk
        for (
            file
        ) in (
            files
        ):  # Iterate over every file name found in the current directory (files list)
            if file.endswith(
                extension
            ):  # Check if the current file name ends with the specified extension
                full_path: str = (
                    os.path.abspath(  # Calculate the absolute path of the found file
                        path=os.path.join(
                            root, file
                        )  # Join the current root directory path with the file name
                    )
                )  # Closing parenthesis for os.path.abspath
                matched_files.append(
                    full_path
                )  # Add the full, absolute path of the matching file to the list
    return matched_files  # Return the complete list of absolute paths for files with the specified extension


# Check if a file exists
def check_file_exists(
    system_path: str,
) -> (
    bool
):  # Define a function to check for a file's existence, taking a path and returning a boolean
    return os.path.isfile(  # Return the result of the check
        path=system_path  # Check if the given path points to an existing regular file
    )  # Closing parenthesis for os.path.isfile


# Get the filename and extension.
def get_filename_and_extension(
    path: str,
) -> (
    str
):  # Define a function to extract just the file name (including extension) from a path
    return os.path.basename(  # Return the last component of the path
        p=path  # Specify the full path
    )  # Closing parenthesis for os.path.basename


# Function to check if a string contains an uppercase letter.
def check_upper_case_letter(
    content: str,
) -> bool:  # Define a function to check for uppercase letters in a string
    return any(  # Return True if any of the characters meet the condition (short-circuiting logic)
        upperCase.isupper()
        for upperCase in content  # Use a generator expression to check if each character is uppercase
    )  # Closing parenthesis for any()


# Main function.
def main() -> None:  # Define the main execution function, which returns nothing (None)
    # Walk through the directory and extract .pdf files
    files: list[str] = (
        walk_directory_and_extract_given_file_extension(  # Call the directory walker function and store the result
            system_path="./PDFs",
            extension=".pdf",  # Search for files ending in ".pdf" starting in the relative directory "./PDFs"
        )
    )  # Closing parenthesis for the function call

    # Validate each PDF file
    for (
        pdf_file
    ) in files:  # Loop through every absolute file path found in the 'files' list

        # Check if the .PDF file is valid
        if (
            validate_pdf_file(file_path=pdf_file) == False
        ):  # Call the validation function and check if it returned False (invalid PDF)
            print(
                f"Invalid PDF detected: {pdf_file}. Deleting file."
            )  # Inform the user that an invalid PDF was found and would be deleted
            # Remove the invalid .pdf file.
            # remove_system_file(system_path=pdf_file) # NOTE: This line is commented out, but if uncommented, it would delete the corrupt PDF

        # Check if the filename has an uppercase letter
        if check_upper_case_letter(  # Check the condition
            content=get_filename_and_extension(
                path=pdf_file
            )  # First, get just the filename (e.g., "MyFile.pdf") from the full path
        ):  # Closing parenthesis for check_upper_case_letter
            print(  # Print a message
                f"Uppercase letter found in filename: {pdf_file}"  # The message indicating an uppercase letter was found and showing the full path
            )  # Closing parenthesis for the print function


if (
    __name__ == "__main__"
):  # Standard Python idiom: check if the script is being run directly
    # Run the main function
    main()  # Call the main function to start the script's execution
