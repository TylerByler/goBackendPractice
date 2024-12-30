const stgOneBtn = document.getElementById("stgOneBtn");
const stgArray = Array.from(document.getElementsByClassName("stg"));
const invoice_number = document.getElementById("invoice_number_temp");
const salesperson_name = document.getElementById("salesperson_name_temp");
const entryInput = document.getElementById("num_entries");

const form = document.getElementById("entryForm");
const addEntry = document.getElementById("addEntry");
const removeEntry = document.getElementById("removeEntry");

var numEntries = 1;

stgOneBtn.addEventListener("click", () => {
    if (invoice_number.value == "" || salesperson_name.value == "") {
        // TELL USER INVALID INPUT
        invoice_number.value = "";
        salesperson_name.value = "";
    } else {
        stgArray.forEach(element => {
            if(element.classList.contains("hidden")) {
                element.classList.remove("hidden");
            } else {
                element.classList.add("hidden");
            }
        });

        document.getElementById("invoice-header").innerText = "Invoice #" + invoice_number.value;
        document.getElementById("salesperson-header").innerText = "Salesperson: " + salesperson_name.value;

        document.getElementById("invoice_number").value = invoice_number.value;
        document.getElementById("salesperson_name").value = salesperson_name.value;
    }
});

addEntry.addEventListener("click", () => {
    numEntries++;
    entryInput.value = numEntries;
    
    const newElement = document.createElement("div");

    newElement.setAttribute("name", `entry${numEntries}`);
    newElement.setAttribute("class", "general-wrapper");

    /*newElement.innerHTML = `
                                <div class="label-wrapper left-end">
                                    <label for="product_number${numEntries}">Product #
                                        <input type="text" name="product_number${numEntries}" id="product_number${numEntries}">
                                    </label>
                                </div>
                                <div class="label-wrapper mid-column">
                                    <label for="product_desc${numEntries}">Product Description
                                        <input type="text" name="product_desc${numEntries}" id="product_desc${numEntries}">
                                    </label>
                                </div>
                                <div class="label-wrapper mid-column">
                                    <label for="color${numEntries}">Color
                                        <input type="text" name="color${numEntries}" id="color${numEntries}">
                                    </label>
                                </div>
                                <div class="label-wrapper mid-column">
                                    <label for="design_number${numEntries}">Design #
                                        <input type="text" name="design_number${numEntries}" id="design_number${numEntries}">
                                    </label>
                                </div>
                                <div class="label-wrapper mid-column">
                                    <label for="font${numEntries}">Font
                                        <input type="text" name="font${numEntries}" id="font${numEntries}">
                                    </label>
                                </div>
                                <div class="label-wrapper right-end">
                                    <label for="engraving_desc${numEntries}">Engraving Description
                                        <input type="text" name="engraving_desc${numEntries}" id="engraving_desc${numEntries}">
                                    </label>
                                </div>
                            `; */

    newElement.innerHTML = `
                                <div class="label-wrapper left-end">
                                    <label for="product_number[]">Product #
                                        <input type="text" name="product_number[]">
                                    </label>
                                </div>
                                <div class="label-wrapper mid-column">
                                    <label for="product_desc[]">Product Description
                                        <input type="text" name="product_desc[]">
                                    </label>
                                </div>
                                <div class="label-wrapper mid-column">
                                    <label for="color[]">Color
                                        <input type="text" name="color[]">
                                    </label>
                                </div>
                                <div class="label-wrapper mid-column">
                                    <label for="design_number[]">Design #
                                        <input type="text" name="design_number[]">
                                    </label>
                                </div>
                                <div class="label-wrapper mid-column">
                                    <label for="font[]">Font
                                        <input type="text" name="font[]">
                                    </label>
                                </div>
                                <div class="label-wrapper right-end">
                                    <label for="engraving_desc[]">Engraving Description
                                        <input type="text" name="engraving_desc[]">
                                    </label>
                                </div>
                            `;

    form.appendChild(newElement);
});

removeEntry.addEventListener("click", () => {
    if (numEntries <= 1) {
        return;
    }

    form.lastChild.remove();
    
    numEntries--;
    entryInput.value = numEntries;
});

printValues = () => {
    console.log("Print New Values")
    console.log(invoice_number.value);
    console.log(salesperson_name.value);
}